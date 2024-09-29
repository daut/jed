package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/daut/simpshop/cmd/api/global"
	"github.com/daut/simpshop/cmd/api/handlers"
	"github.com/daut/simpshop/db"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "user=daut dbname=simpshop sslmode=verify-full")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &global.Application{
		Queries:  db.New(conn),
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	h := &handlers.Handler{App: app}

	router := http.NewServeMux()
	router.HandleFunc("POST /products", h.ProductCreate)
	router.HandleFunc("GET /products", h.ProductList)
	router.HandleFunc("GET /products/{id}", h.ProductRead)
	router.HandleFunc("DELETE /products/{id}", h.ProductDelete)

	infoLog.Printf("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}
