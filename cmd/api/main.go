package main

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "user=daut dbname=simpshop sslmode=verify-full")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	app := NewApp(conn)

	app.Logger.Info.Printf("Starting server on :8080")
	http.ListenAndServe(":8080", app.Router)
}
