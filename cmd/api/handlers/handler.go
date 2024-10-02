package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/daut/simpshop/db"
	"github.com/daut/simpshop/internal/utils"
)

type Handler struct {
	Queries *db.Queries
	Logger  *utils.Logger
}

func New(queries *db.Queries, logger *utils.Logger) *Handler {
	return &Handler{
		Queries: queries,
		Logger:  logger,
	}
}

// Returns a 500 Internal Server Error response to the client
func (h *Handler) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.Logger.Error.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Returns a 400 Bad Request response to the client
func (h *Handler) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Returns a 404 Not Found response to the client
func (h *Handler) NotFound(w http.ResponseWriter) {
	h.ClientError(w, http.StatusNotFound)
}

// Writes a JSON response to the client
func (h *Handler) WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) {
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
