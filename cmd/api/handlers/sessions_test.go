package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/daut/jed/internal/assert"
	"github.com/daut/jed/internal/testutils"
	"github.com/daut/jed/sqlc"
)

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
}

func TestSessionCreate(t *testing.T) {
	t.Parallel()
	queries := []string{"insert into admins (username, password) values ('admin', crypt('password', gen_salt('bf')));"}
	dbr := testutils.NewDBResources(t, queries)
	defer dbr.Close(t)
	handlers := initHandlers(dbr.Pool)

	tests := []struct {
		Name           string
		Body           string
		ExpectedStatus int
	}{
		{Name: "Valid credentials", Body: `{"username":"admin","password":"password"}`, ExpectedStatus: http.StatusCreated},
		{Name: "Invalid credentials", Body: `{"username":"admin","password":"invalid"}`, ExpectedStatus: http.StatusUnauthorized},
		{Name: "Invalid username", Body: `{"username":"invalid","password":"password"}`, ExpectedStatus: http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.Body))
			w := httptest.NewRecorder()
			handlers.SessionCreate(w, req)
			resp := w.Result()
			var response LoginResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			assert.Nil(t, err)
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
			assert.NotNil(t, response.Token)
			assert.NotNil(t, response.ExpiresAt)
		})
	}
}

func TestSessionDelete(t *testing.T) {
	t.Parallel()
	queries := []string{"insert into admins (username, password) values ('admin', crypt('password', gen_salt('bf')));"}
	dbr := testutils.NewDBResources(t, queries)
	defer dbr.Close(t)
	handlers := initHandlers(dbr.Pool)
	tests := []struct {
		Name           string
		LoginBody      string
		ID             string
		ExpectedStatus int
	}{
		{Name: "Invalid credentials", LoginBody: `{"username":"admin","password":"password"}`, ID: "invalid", ExpectedStatus: http.StatusBadRequest},
		{Name: "Valid credentials", LoginBody: `{"username":"admin","password":"password"}`, ID: "1", ExpectedStatus: http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.LoginBody))
			w := httptest.NewRecorder()
			handlers.SessionCreate(w, req)
			resp := w.Result()
			var response LoginResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			assert.Nil(t, err)
			req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/sessions/%s", tt.ID), nil)
			w = httptest.NewRecorder()
			req.SetPathValue("id", tt.ID)
			id, err := strconv.Atoi(tt.ID)
			if err == nil {
				req = req.WithContext(context.WithValue(req.Context(), "admin", sqlc.Admin{ID: int32(id)}))
			}
			handlers.SessionDelete(w, req)
			resp = w.Result()
			err = json.NewDecoder(resp.Body).Decode(&response)
			assert.Nil(t, err)
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
		})
	}
}
