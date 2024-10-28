package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/daut/jed/internal/assert"
	"github.com/daut/jed/internal/testutils"
)

func TestLogin(t *testing.T) {
	t.Parallel()
	queries := []string{"insert into admins (username, password) values ('admin', crypt('password', gen_salt('bf')));"}
	dbr := testutils.NewDBResources(t, queries)
	defer dbr.Close(t)
	handlers := initHandlers(dbr.Conn)

	tests := []struct {
		Name           string
		Body           string
		ExpectedStatus int
	}{
		{Name: "Valid credentials", Body: `{"username":"admin","password":"password"}`, ExpectedStatus: 201},
		{Name: "Invalid credentials", Body: `{"username":"admin","password":"invalid"}`, ExpectedStatus: 401},
		{Name: "Invalid username", Body: `{"username":"invalid","password":"password"}`, ExpectedStatus: 401},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.Body))
			w := httptest.NewRecorder()
			handlers.Login(w, req)
			resp := w.Result()
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
		})
	}
}
