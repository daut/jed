package handlers

import (
	"github.com/daut/jed/cmd/api/helpers"
	"github.com/daut/jed/internal/utils"
	db "github.com/daut/jed/sqlc"
	"github.com/jackc/pgx/v5"
)

func initHandlers(conn *pgx.Conn) *Handler {
	queries := db.New(conn)
	logger := utils.NewLogger()
	responseHelper := helpers.NewResponse(logger)
	handlers := New(queries, logger, responseHelper)
	return handlers
}
