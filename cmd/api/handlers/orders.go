package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/daut/jed/cmd/api/consts"
	"github.com/daut/jed/cmd/api/dto"
	db "github.com/daut/jed/sqlc"
)

func (handler *Handler) OrderCreate(w http.ResponseWriter, r *http.Request) {
	var input dto.OrderCreateRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		handler.Response.ClientError(w, consts.ErrInvalidInput, http.StatusBadRequest)
	}

	args := &db.CreateOrderParams{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Phone:     input.Phone,
		Address:   input.Address,
		City:      input.City,
	}
	order, err := handler.Queries.CreateOrder(r.Context(), *args)
	if err != nil {
		handler.Response.ServerError(w, err)
	}

	handler.Response.WriteJSON(w, http.StatusCreated, order, nil)
}
