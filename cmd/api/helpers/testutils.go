package helpers

import (
	"github.com/daut/jed/internal/utils"
	db "github.com/daut/jed/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TestResources struct {
	Queries  *db.Queries
	Logger   *utils.Logger
	Response *Response
}

func NewTestResources(conn *pgxpool.Pool) *TestResources {
	queries := db.New(conn)
	logger := utils.NewLogger()
	responseHelper := NewResponse(logger)
	return &TestResources{
		Queries:  queries,
		Logger:   logger,
		Response: responseHelper,
	}
}
