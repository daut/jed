package router

import (
	"net/http"

	"github.com/daut/simpshop/cmd/api/handlers"
	"github.com/daut/simpshop/db"
	"github.com/daut/simpshop/internal/utils"
)

func New(queries *db.Queries, logger *utils.Logger) http.Handler {
	handlers := handlers.New(queries, logger)
	router := http.NewServeMux()
	router.HandleFunc("POST /products", handlers.ProductCreate)
	router.HandleFunc("GET /products", handlers.ProductList)
	router.HandleFunc("GET /products/{id}", handlers.ProductRead)
	router.HandleFunc("PUT /products/{id}", handlers.ProductUpdate)
	router.HandleFunc("DELETE /products/{id}", handlers.ProductDelete)
	return router
}
