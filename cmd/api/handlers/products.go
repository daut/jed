package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/daut/jed/internal/utils"
	db "github.com/daut/jed/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func (handler *Handler) ProductCreate(w http.ResponseWriter, r *http.Request) {
	// TODO: Needs admin authentication

	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		handler.Response.ClientError(w, "invalid input", http.StatusBadRequest)
		return
	}

	price, err := utils.ConvertToPGNumeric(input.Price)
	if err != nil {
		handler.Response.ClientError(w, "invalid price", http.StatusBadRequest)
		return
	}
	args := &db.CreateProductParams{
		Name:        input.Name,
		Description: pgtype.Text{String: input.Description, Valid: true},
		Price:       *price,
	}
	prod, err := handler.Queries.CreateProduct(r.Context(), *args)
	if err != nil {
		handler.Response.ServerError(w, err)
		return
	}

	handler.Response.WriteJSON(w, http.StatusCreated, prod, nil)
}

func (handler *Handler) ProductRead(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		handler.Response.ClientError(w, "invalid id", http.StatusBadRequest)
		return
	}
	prod, err := handler.Queries.GetProduct(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handler.Response.NotFound(w)
		} else {
			handler.Response.ServerError(w, err)
		}
		return
	}

	handler.Response.WriteJSON(w, http.StatusOK, prod, nil)
}

func (handler *Handler) ProductList(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		handler.Response.ClientError(w, "invalid page", http.StatusBadRequest)
		return
	}

	args := &db.GetProductsParams{
		Limit:  10,
		Offset: (int32(page) - 1) * 10,
	}
	products, err := handler.Queries.GetProducts(r.Context(), *args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handler.Response.NotFound(w)
		} else {
			handler.Response.ServerError(w, err)
		}
		return
	}

	if len(products) == 0 {
		handler.Response.NotFound(w)
		return
	}

	handler.Response.WriteJSON(w, http.StatusOK, products, nil)
}

func (handler *Handler) ProductUpdate(w http.ResponseWriter, r *http.Request) {
	// TODO: Needs admin authentication

	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		handler.Response.ClientError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input struct {
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Price       *float64 `json:"price"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		handler.Response.ClientError(w, "invalid input", http.StatusBadRequest)
		return
	}

	if input.Name == nil || input.Description == nil || input.Price == nil {
		handler.Response.ClientError(w, "invalid input", http.StatusBadRequest)
		return
	}

	price, err := utils.ConvertToPGNumeric(*input.Price)
	if err != nil {
		handler.Response.ClientError(w, "invalid price", http.StatusBadRequest)
		return
	}
	args := &db.UpdateProductParams{
		Name:        *input.Name,
		Description: pgtype.Text{String: *input.Description, Valid: true},
		Price:       *price,
		ID:          int32(id),
	}
	product, err := handler.Queries.UpdateProduct(r.Context(), *args)
	if err != nil {
		handler.Response.ServerError(w, err)
		return
	}

	handler.Response.WriteJSON(w, http.StatusOK, product, nil)
}

func (handler *Handler) ProductDelete(w http.ResponseWriter, r *http.Request) {
	// TODO: Needs admin authentication

	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		handler.Response.ClientError(w, "invalid id", http.StatusBadRequest)
		return
	}

	prod, err := handler.Queries.DeleteProduct(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handler.Response.NotFound(w)
		} else {
			handler.Response.ServerError(w, err)
		}
		return
	}

	handler.Logger.Info.Printf("Product %v deleted", prod.Name)

	handler.Response.WriteJSON(w, http.StatusNoContent, nil, nil)
}
