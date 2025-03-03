package handlers

import (
	"github.com/daut/jed/cmd/api/helpers"
	"github.com/daut/jed/internal/utils"
	db "github.com/daut/jed/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initHandlers(pool *pgxpool.Pool) *Handler {
	queries := db.New(pool)
	logger := utils.NewLogger()
	responseHelper := helpers.NewResponse(logger)
	handlers := New(queries, logger, responseHelper, pool)
	return handlers
}
