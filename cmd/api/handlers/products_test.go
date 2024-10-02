package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daut/simpshop/db"
	"github.com/daut/simpshop/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
)

var user = "test"
var password = "test"
var databaseName = "test_shop"

func TestProductRead(t *testing.T) {
	p := postgres.Preset(
		postgres.WithUser(user, password),
		postgres.WithDatabase(databaseName),
		postgres.WithQueries("insert into products (name, description, price) values ('test', 'test', 1.00)"),
		postgres.WithQueriesFile("../../../sqlc/schema.sql"),
		postgres.WithTimezone("Europe/Belgrade"),
	)

	container, err := gnomock.Start(p)
	if err != nil {
		t.Fatal(err)
	}
	defer gnomock.Stop(container)

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "user=daut dbname=simpshop sslmode=verify-full")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	handlers := New(db.New(conn), utils.NewLogger())

	req := httptest.NewRequest("GET", "/products/1", nil)
	w := httptest.NewRecorder()
	handlers.ProductRead(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}
