package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daut/simpshop/db"
	"github.com/daut/simpshop/internal/assert"
	"github.com/daut/simpshop/internal/testutils"
	"github.com/daut/simpshop/internal/utils"
	"github.com/orlangure/gnomock"
)

func TestProductRead(t *testing.T) {
	queries := []string{"insert into products (name, description, price) values ('product1', 'good product', 100)"}
	container := testutils.NewDBContainer(t, queries)
	defer gnomock.Stop(container)

	conn := testutils.NewDBConn(t, container)
	defer conn.Close(context.Background())

	handlers := New(db.New(conn), utils.NewLogger())

	tests := []struct {
		Name           string
		ID             string
		Expected       string
		Actual         string
		ExpectedStatus int
		ActualStatus   int
	}{
		{Name: "Product exists", ID: "1", Expected: `{"id":1,"name":"product1","description":"good product","price":100}`, ExpectedStatus: http.StatusOK},
		{Name: "Product does not exist", ID: "2", Expected: `{"error":"product not found"}`, ExpectedStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/products/"+tt.ID, nil)
			w := httptest.NewRecorder()
			req.SetPathValue("id", tt.ID)
			handlers.ProductRead(w, req)
			resp := w.Result()
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
			// assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}
}
