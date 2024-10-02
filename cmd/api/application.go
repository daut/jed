package main

import (
	"net/http"

	"github.com/daut/simpshop/cmd/api/router"
	"github.com/daut/simpshop/db"
	"github.com/daut/simpshop/internal/utils"
	"github.com/jackc/pgx/v5"
)

type Application struct {
	Queries *db.Queries
	Logger  *utils.Logger
	Router  http.Handler
}

func New(conn *pgx.Conn) *Application {
	queries := db.New(conn)
	logger := utils.NewLogger()
	router := router.New(queries, logger)
	return &Application{
		Queries: queries,
		Logger:  logger,
		Router:  router,
	}
}
