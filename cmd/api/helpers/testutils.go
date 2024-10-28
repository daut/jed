package helpers

import (
	"github.com/daut/jed/internal/utils"
	db "github.com/daut/jed/sqlc"
	"github.com/jackc/pgx/v5"
)

type TestResources struct {
	Queries  *db.Queries
	Logger   *utils.Logger
	Response *Response
}

func NewTestResources(conn *pgx.Conn) *TestResources {
	queries := db.New(conn)
	logger := utils.NewLogger()
	responseHelper := NewResponse(logger)
	return &TestResources{
		Queries:  queries,
		Logger:   logger,
		Response: responseHelper,
	}
}
