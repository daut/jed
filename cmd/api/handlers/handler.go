package handlers

import (
	"github.com/daut/jed/cmd/api/helpers"
	"github.com/daut/jed/internal/utils"
	db "github.com/daut/jed/sqlc"
)

type Handler struct {
	Queries  *db.Queries
	Logger   *utils.Logger
	Response *helpers.Response
}

func New(queries *db.Queries, logger *utils.Logger, response *helpers.Response) *Handler {
	return &Handler{
		Queries:  queries,
		Logger:   logger,
		Response: response,
	}
}
