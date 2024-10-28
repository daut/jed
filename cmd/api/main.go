package main

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	conn, err := pgxpool.New(context.Background(), "user=daut dbname=jed sslmode=verify-full")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	app := NewApp(conn)

	app.Logger.Info.Printf("Starting server on :8080")
	http.ListenAndServe(":8080", app.Router)
}
