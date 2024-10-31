package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/daut/jed/internal/assert"
	"github.com/daut/jed/internal/testutils"
)

func TestProductRead(t *testing.T) {
	t.Parallel()
	queries := []string{"insert into products (name, description, price) values ('product1', 'good product', 100)"}
	dbr := testutils.NewDBResources(t, queries)
	defer dbr.Close(t)
	handlers := initHandlers(dbr.Pool)

	tests := []struct {
		Name           string
		ID             string
		Expected       string
		ExpectedStatus int
	}{
		{Name: "Product exists", ID: "1", Expected: `{"id":1,"name":"product1","description":"good product","price":100.00}`, ExpectedStatus: http.StatusOK},
		{Name: "Product does not exist", ID: "2", Expected: `{"message":"the requested resource could not be found"}`, ExpectedStatus: http.StatusNotFound},
		{Name: "Invalid ID", ID: "invalid", Expected: `{"message":"invalid id"}`, ExpectedStatus: http.StatusBadRequest},
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
	dbr := testutils.NewDBResources(t, queries)
	defer dbr.Close(t)
	handlers := initHandlers(dbr.Pool)

	tests := []struct {
		Name           string
		Offset         string
		Expected       string
		ExpectedStatus int
	}{
		{Name: "List products", Offset: "1", Expected: `[{"id":1,"name":"product1","description":"good product","price":100.00}]`, ExpectedStatus: http.StatusOK},
		{Name: "No offset", Offset: "", Expected: `[{"id":1,"name":"product1","description":"good product","price":100.00}]`, ExpectedStatus: http.StatusOK},
		{Name: "Invalid offset", Offset: "invalid", Expected: `{"message":"invalid page"}`, ExpectedStatus: http.StatusBadRequest},
		{Name: "No products", Offset: "5", Expected: `{"message":"the requested resource could not be found"}`, ExpectedStatus: http.StatusNotFound},
		{Name: "Negative offset", Offset: "-1", Expected: `{"message":"invalid page"}`, ExpectedStatus: http.StatusBadRequest},
		{Name: "Zero offset", Offset: "0", Expected: `{"message":"invalid page"}`, ExpectedStatus: http.StatusBadRequest},
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

func TestProductCreate(t *testing.T) {
	t.Parallel()
	queries := []string{}
	dbr := testutils.NewDBResources(t, queries)
	defer dbr.Close(t)
	handlers := initHandlers(dbr.Pool)

	tests := []struct {
		Name           string
		Body           string
		Expected       string
		ExpectedStatus int
	}{
		{Name: "Valid product", Body: `{"name":"product1","description":"good product","price":100}`, Expected: `{"id":1,"name":"product1","description":"good product","price":100.00}`, ExpectedStatus: http.StatusCreated},
		{Name: "Invalid JSON", Body: `{"name":"product1"`, Expected: `{"message":"invalid input"}`, ExpectedStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/products", strings.NewReader(tt.Body))
			w := httptest.NewRecorder()
			handlers.ProductCreate(w, req)
			resp := w.Result()
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
			assert.Equal(t, tt.Expected, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestProductUpdate(t *testing.T) {
	t.Parallel()
	queries := []string{"insert into products (name, description, price) values ('product1', 'good product', 1000)"}
	dbr := testutils.NewDBResources(t, queries)
	defer dbr.Close(t)
	handlers := initHandlers(dbr.Pool)

	tests := []struct {
		Name           string
		ID             string
		Body           string
		Expected       string
		ExpectedStatus int
	}{
		{Name: "Update price", ID: "1", Body: `{"name":"product1","description":"good product","price":1000}`, Expected: `{"id":1,"name":"product1","description":"good product","price":1000.00}`, ExpectedStatus: http.StatusOK},
		{Name: "Invalid ID", ID: "invalid", Body: `{"name":"product1","description":"good product","price":1000}`, Expected: `{"message":"invalid id"}`, ExpectedStatus: http.StatusBadRequest},
		{Name: "Missing price", ID: "1", Body: `{"name":"product1","description":"good product"}`, Expected: `{"message":"invalid input"}`, ExpectedStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", "/products/"+tt.ID, strings.NewReader(tt.Body))
			w := httptest.NewRecorder()
			req.SetPathValue("id", tt.ID)
			handlers.ProductUpdate(w, req)
			resp := w.Result()
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
			assert.Equal(t, tt.Expected, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestProductDelete(t *testing.T) {
	t.Parallel()
	queries := []string{"insert into products (name, description, price) values ('product1', 'good product', 1000)"}
	dbr := testutils.NewDBResources(t, queries)
	defer dbr.Close(t)
	handlers := initHandlers(dbr.Pool)

	tests := []struct {
		Name           string
		ID             string
		Expected       string
		ExpectedStatus int
	}{
		{Name: "Delete product", ID: "1", Expected: "null", ExpectedStatus: http.StatusNoContent},
		{Name: "Product does not exist", ID: "2", Expected: `{"message":"the requested resource could not be found"}`, ExpectedStatus: http.StatusNotFound},
		{Name: "Invalid ID", ID: "invalid", Expected: `{"message":"invalid id"}`, ExpectedStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/products/"+tt.ID, nil)
			w := httptest.NewRecorder()
			req.SetPathValue("id", tt.ID)
			handlers.ProductDelete(w, req)
			resp := w.Result()
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
			assert.Equal(t, tt.Expected, strings.TrimSpace(w.Body.String()))
		})
	}
}
