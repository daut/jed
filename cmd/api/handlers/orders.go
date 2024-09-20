package handlers

import (
	"fmt"
	"net/http"
)

func OrderCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Order Create"))
}

func OrderRead(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Print(id)
	w.Write([]byte("Order Read"))
}

func OrderList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Order List"))
}

func OrderUpdate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Order Update"))
}

func OrderDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Print(id)
	w.Write([]byte("Order Delete"))
}
