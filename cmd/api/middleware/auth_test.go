package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/daut/jed/cmd/api/helpers"
	"github.com/daut/jed/internal/assert"
	"github.com/daut/jed/internal/testutils"
	"github.com/daut/jed/internal/tokens"
	"github.com/daut/jed/internal/utils"
	"github.com/daut/jed/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestAuth(t *testing.T) {
	t.Parallel()
	token, _ := tokens.GenerateToken(1, time.Hour)
	queries := []string{
		"insert into admins (username, password) values ('admin', 'password')",
	}
	dbr := testutils.NewDBResources(t, queries)
	defer dbr.Close(t)
	res := helpers.NewTestResources(dbr.Pool)
	middleware := New(res.Queries, res.Logger, res.Response)
	middleware.Queries.SaveToken(context.Background(), sqlc.SaveTokenParams{
		Hash:      token.Hash,
		AdminID:   token.AdminID,
		ExpiresAt: pgtype.Timestamptz{Time: token.ExpiresAt, Valid: true},
	})
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	tests := []struct {
		Name           string
		Authorization  *string
		ExpectedStatus int
	}{
		{Name: "Missing authorization header", Authorization: nil, ExpectedStatus: http.StatusBadRequest},
		{Name: "Missing token", Authorization: utils.StrPtr("Bearer"), ExpectedStatus: http.StatusBadRequest},
		{Name: "Invalid token", Authorization: utils.StrPtr("Bearer random"), ExpectedStatus: http.StatusUnauthorized},
		{Name: "Valid token", Authorization: utils.StrPtr("Bearer " + token.PlainText), ExpectedStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.Authorization != nil {
				req.Header.Set("Authorization", *tt.Authorization)
			}
			w := httptest.NewRecorder()
			middleware.Auth(next).ServeHTTP(w, req)
			resp := w.Result()
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
		})
	}
}
