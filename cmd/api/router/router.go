package router

import (
	"net/http"

	"github.com/daut/jed/cmd/api/handlers"
	"github.com/daut/jed/cmd/api/helpers"
	"github.com/daut/jed/internal/utils"
	db "github.com/daut/jed/sqlc"
)

func New(queries *db.Queries, logger *utils.Logger) http.Handler {
	responseHelper := helpers.NewResponse(logger)
	handlers := handlers.New(queries, logger, responseHelper)
	router := http.NewServeMux()
	router.HandleFunc("POST /products", handlers.ProductCreate)
	router.HandleFunc("GET /products", handlers.ProductList)
	router.HandleFunc("GET /products/{id}", handlers.ProductRead)
	router.HandleFunc("PUT /products/{id}", handlers.ProductUpdate)
	router.HandleFunc("DELETE /products/{id}", handlers.ProductDelete)

	router.HandleFunc("GET /admins/{id}", handlers.AdminRead)

	router.HandleFunc("POST /login", handlers.Login)
	return router
}
