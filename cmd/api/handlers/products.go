package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/daut/simpshop/cmd/api/global"
	"github.com/daut/simpshop/db"
)

type Handler struct {
	App *global.Application
}

func ProductCreate(w http.ResponseWriter, r *http.Request) {
	// Needs admin authentication
	w.Write([]byte("Product Create"))
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

func ProductDelete(w http.ResponseWriter, r *http.Request) {
	// Needs admin authentication
	id := r.PathValue("id")
	fmt.Print(id)
	w.Write([]byte("Product Delete"))
}
