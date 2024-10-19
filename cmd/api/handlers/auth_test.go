package handlers

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/daut/jed/internal/assert"
	"github.com/daut/jed/internal/testutils"
	"github.com/daut/jed/internal/utils"
	db "github.com/daut/jed/sqlc"
	"github.com/orlangure/gnomock"
)

func TestLogin(t *testing.T) {
	t.Parallel()
	queries := []string{"insert into admins (username, password) values ('admin', crypt('password', gen_salt('bf')));"}
	container := testutils.NewDBContainer(t, queries)
	defer gnomock.Stop(container)

	conn := testutils.NewDBConn(t, container)
	defer conn.Close(context.Background())

	handlers := New(db.New(conn), utils.NewLogger())

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
