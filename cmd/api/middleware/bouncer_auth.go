package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daut/jed/cmd/api/helpers"
	"github.com/daut/jed/internal/assert"
	"github.com/daut/jed/internal/testutils"
)

func TestRequireAdminUser(t *testing.T) {
	t.Parallel()
	queries := []string{}
	dbr := testutils.NewDBResources(t, queries)
	defer dbr.Close(t)
	res := helpers.NewTestResources(dbr.Pool)
	middleware := New(res.Queries, res.Logger, res.Response)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	tests := []struct {
		Name           string
		Admin          bool
		ExpectedStatus int
	}{
		{Name: "Admin user", Admin: true, ExpectedStatus: http.StatusOK},
		{Name: "Non-admin user", Admin: false, ExpectedStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.Admin {
				req = req.WithContext(context.WithValue(req.Context(), "admin", struct{}{}))
			}
			w := httptest.NewRecorder()
			middleware.RequireAdminUser(next).ServeHTTP(w, req)
			resp := w.Result()
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
		})
	}
}
