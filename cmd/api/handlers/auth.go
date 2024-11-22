package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/daut/jed/internal/consts"
	"github.com/daut/jed/internal/tokens"
	"github.com/daut/jed/internal/validator"
	db "github.com/daut/jed/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		handler.Response.ClientError(w, consts.ErrInvalidInput, http.StatusBadRequest)
		return
	}

	v := validator.New()
	v.IsNotEmpty(input.Username, "username", consts.ErrMissingField)
	v.IsNotEmpty(input.Password, "password", consts.ErrMissingField)
	if v.HasErrors() {
		handler.Response.FailedValidation(w, v.Errors)
	}

	admin, err := handler.Queries.GetAdmin(r.Context(), input.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handler.Response.ClientError(w, consts.ErrInvalidCreds, http.StatusUnauthorized)
		} else {
			handler.Response.ServerError(w, err)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password))
	if err != nil {
		handler.Response.ClientError(w, consts.ErrInvalidCreds, http.StatusUnauthorized)
		return
	}

	token, err := tokens.GenerateToken(admin.ID, 72*time.Hour)
	if err != nil {
		handler.Response.ServerError(w, err)
		return
	}

	args := &db.SaveTokenParams{
		Hash:      token.Hash,
		AdminID:   token.AdminID,
		ExpiresAt: pgtype.Timestamptz{Time: token.ExpiresAt, Valid: true},
	}
	_, err = handler.Queries.SaveToken(r.Context(), *args)
	if err != nil {
		handler.Response.ServerError(w, err)
		return
	}

	handler.Response.WriteJSON(w, http.StatusCreated, token, nil)
}

func (handler *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: implement logout
	handler.Response.WriteJSON(w, http.StatusOK, nil, nil)
}
