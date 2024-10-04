package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/daut/simpshop/db"
	"github.com/daut/simpshop/internal/assert"
	"github.com/daut/simpshop/internal/testutils"
	"github.com/daut/simpshop/internal/utils"
	"github.com/orlangure/gnomock"
)

func TestProductRead(t *testing.T) {
	t.Parallel()
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
		ExpectedStatus int
	}{
		{Name: "Product exists", ID: "1", Expected: `{"id":1,"name":"product1","description":"good product","price":100.00}`, ExpectedStatus: http.StatusOK},
		{Name: "Product does not exist", ID: "2", Expected: "Not Found", ExpectedStatus: http.StatusNotFound},
		{Name: "Invalid ID", ID: "invalid", Expected: "Bad Request", ExpectedStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/products/"+tt.ID, nil)
			w := httptest.NewRecorder()
			req.SetPathValue("id", tt.ID)
			handlers.ProductRead(w, req)
			resp := w.Result()
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
			assert.Equal(t, tt.Expected, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestProductList(t *testing.T) {
	t.Parallel()
	queries := []string{"insert into products (name, description, price) values ('product1', 'good product', 100)"}
	container := testutils.NewDBContainer(t, queries)
	defer gnomock.Stop(container)

	conn := testutils.NewDBConn(t, container)
	defer conn.Close(context.Background())

	handlers := New(db.New(conn), utils.NewLogger())

	tests := []struct {
		Name           string
		Offset         string
		Expected       string
		ExpectedStatus int
	}{
		{Name: "List products", Offset: "1", Expected: `[{"id":1,"name":"product1","description":"good product","price":100.00}]`, ExpectedStatus: http.StatusOK},
		{Name: "No offset", Offset: "", Expected: `[{"id":1,"name":"product1","description":"good product","price":100.00}]`, ExpectedStatus: http.StatusOK},
		{Name: "Invalid offset", Offset: "invalid", Expected: "Bad Request", ExpectedStatus: http.StatusBadRequest},
		{Name: "No products", Offset: "5", Expected: "Not Found", ExpectedStatus: http.StatusNotFound},
		{Name: "Negative offset", Offset: "-1", Expected: "Bad Request", ExpectedStatus: http.StatusBadRequest},
		{Name: "Zero offset", Offset: "0", Expected: "Bad Request", ExpectedStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/products?page="+tt.Offset, nil)
			w := httptest.NewRecorder()
			req.URL.Query().Set("page", tt.Offset)
			handlers.ProductList(w, req)
			resp := w.Result()
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
			assert.Equal(t, tt.Expected, strings.TrimSpace(w.Body.String()))
		})
	}
}
