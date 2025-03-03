package handlers

import (
	"github.com/daut/jed/cmd/api/helpers"
	"github.com/daut/jed/internal/utils"
	db "github.com/daut/jed/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	Queries  *db.Queries
	Logger   *utils.Logger
	Response *helpers.Response
	Pool     *pgxpool.Pool
}

func New(queries *db.Queries, logger *utils.Logger, response *helpers.Response, pool *pgxpool.Pool) *Handler {
	return &Handler{
		Queries:  queries,
		Logger:   logger,
		Response: response,
		Pool:     pool,
	}
}
