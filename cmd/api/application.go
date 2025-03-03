package main

import (
	"net/http"

	"github.com/daut/jed/cmd/api/router"
	"github.com/daut/jed/internal/utils"
	db "github.com/daut/jed/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Logger *utils.Logger
	Router http.Handler
}

func NewApp(pool *pgxpool.Pool) *Application {
	queries := db.New(pool)
	logger := utils.NewLogger()
	router := router.New(queries, logger, pool)
	return &Application{
		Logger: logger,
		Router: router,
	}
}
