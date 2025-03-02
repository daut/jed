package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/daut/jed/internal/assert"
	"github.com/daut/jed/internal/testutils"
)

func TestOrderCreate(t *testing.T) {
	t.Parallel()
	dbr := testutils.NewDBResources(t, []string{})
	defer dbr.Close(t)
	handlers := initHandlers(dbr.Pool)

	tests := []struct {
		Name           string
		Body           string
		ExpectedStatus int
	}{}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/orders", strings.NewReader(tt.Body))
			w := httptest.NewRecorder()
			handlers.OrderCreate(w, req)
			resp := w.Result()
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
		})
	}
}
