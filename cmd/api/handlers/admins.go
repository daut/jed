package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/daut/jed/cmd/api/consts"
	"github.com/daut/jed/internal/validator"
	db "github.com/daut/jed/sqlc"
)

type Admin struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) AdminRead(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	v := validator.New()
	v.IsNotEmpty(username, "username", consts.ErrMissingField)
	if v.HasErrors() {
		h.Response.FailedValidation(w, v.Errors)
	}

	admin, err := h.Queries.GetAdmin(r.Context(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Response.NotFound(w)
		} else {
			h.Response.ServerError(w, err)
		}
		return
	}

	response := &Admin{
		ID:       admin.ID,
		Username: admin.Username,
	}

	h.Response.WriteJSON(w, http.StatusOK, response, nil)
}

func (h *Handler) AdminList(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		h.Response.ClientError(w, "invalid page", http.StatusBadRequest)
		return
	}

	args := &db.ListAdminsParams{
		Limit:  10,
		Offset: (int32(page) - 1) * 10,
	}
	admins, err := h.Queries.ListAdmins(r.Context(), *args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Response.NotFound(w)
		} else {
			h.Response.ServerError(w, err)
		}
		return
	}

	if len(admins) == 0 {
		h.Response.NotFound(w)
		return
	}

	h.Response.WriteJSON(w, http.StatusOK, admins, nil)
}
