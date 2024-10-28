package middleware

import (
	"github.com/daut/jed/cmd/api/helpers"
	"github.com/daut/jed/internal/utils"
	db "github.com/daut/jed/sqlc"
)

type Middleware struct {
	Response *helpers.Response
	Logger   *utils.Logger
	Queries  *db.Queries
}

func New(queries *db.Queries, logger *utils.Logger, response *helpers.Response) *Middleware {
	return &Middleware{
		Response: response,
		Logger:   logger,
		Queries:  queries,
	}
}
