package router

import (
	"net/http"

	"github.com/daut/jed/cmd/api/handlers"
	"github.com/daut/jed/cmd/api/helpers"
	"github.com/daut/jed/cmd/api/middleware"
	"github.com/daut/jed/internal/utils"
	db "github.com/daut/jed/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/justinas/alice"
)

func New(queries *db.Queries, logger *utils.Logger, pool *pgxpool.Pool) http.Handler {
	responseHelper := helpers.NewResponse(logger)
	handlers := handlers.New(queries, logger, responseHelper, pool)
	router := http.NewServeMux()

	mw := middleware.New(queries, logger, responseHelper)

	isAdmin := alice.New(mw.Auth, mw.RequireAdminUser)

	// public
	router.HandleFunc("GET /products", handlers.ProductList)
	router.HandleFunc("GET /products/{id}", handlers.ProductRead)

	router.HandleFunc("POST /sessions", handlers.SessionCreate)

	// admin
	router.Handle("POST /products", isAdmin.ThenFunc(handlers.ProductCreate))
	router.Handle("PUT /products/{id}", isAdmin.ThenFunc(handlers.ProductUpdate))
	router.Handle("DELETE /products/{id}", isAdmin.ThenFunc(handlers.ProductDelete))

	router.Handle("GET /admins/{username}", isAdmin.ThenFunc(handlers.AdminRead))
	router.Handle("GET /admins", isAdmin.ThenFunc(handlers.AdminList))

	router.Handle("DELETE /sessions/{id}", isAdmin.ThenFunc(handlers.SessionDelete))

	return router
}
