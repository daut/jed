package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/daut/simpshop/cmd/api/global"
	"github.com/daut/simpshop/db"
	"github.com/daut/simpshop/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	App *global.Application
}

func (h *Handler) ProductCreate(w http.ResponseWriter, r *http.Request) {
	// Needs admin authentication

	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.App.ClientError(w, http.StatusBadRequest)
		return
	}

	price, err := utils.ConvertToPGNumeric(input.Price)
	if err != nil {
		h.App.ClientError(w, http.StatusBadRequest)
		return
	}
	args := &db.CreateProductParams{
		Name:        pgtype.Text{String: input.Name, Valid: true},
		Description: pgtype.Text{String: input.Description, Valid: true},
		Price:       *price,
	}
	prod, err := h.App.Queries.CreateProduct(r.Context(), *args)
	if err != nil {
		h.App.ServerError(w, err)
		return
	}

	h.App.WriteJSON(w, http.StatusCreated, prod, nil)
}

func (h *Handler) ProductRead(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.App.ClientError(w, http.StatusBadRequest)
		return
	}
	prod, err := h.App.Queries.GetProduct(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.App.NotFound(w)
		} else {
			h.App.ServerError(w, err)
		}
		return
	}

	h.App.WriteJSON(w, http.StatusOK, prod, nil)
}

func (h *Handler) ProductList(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1"
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		h.App.ClientError(w, http.StatusBadRequest)
		return
	}

	args := &db.GetProductsParams{
		Limit:  10,
		Offset: (int32(page) - 1) * 10,
	}
	products, err := h.App.Queries.GetProducts(r.Context(), *args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.App.NotFound(w)
		} else {
			h.App.ServerError(w, err)
		}
		return
	}

	if len(products) == 0 {
		h.App.NotFound(w)
		return
	}

	h.App.WriteJSON(w, http.StatusOK, products, nil)
}

func ProductUpdate(w http.ResponseWriter, r *http.Request) {
	// Needs admin authentication
	w.Write([]byte("Product Update"))
}

func (h *Handler) ProductDelete(w http.ResponseWriter, r *http.Request) {
	// Needs admin authentication

	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.App.ClientError(w, http.StatusBadRequest)
	}

	_, err = h.App.Queries.DeleteProduct(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.App.NotFound(w)
		} else {
			h.App.ServerError(w, err)
		}
	}

	h.App.InfoLog.Printf("Product %d deleted", id)

	h.App.WriteJSON(w, http.StatusNoContent, nil, nil)
}
