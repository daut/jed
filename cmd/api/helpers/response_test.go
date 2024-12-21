package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/daut/jed/internal/assert"
	"github.com/daut/jed/internal/consts"
	"github.com/daut/jed/internal/utils"
)

func TestClientError(t *testing.T) {
	t.Parallel()
	responseHelper := NewResponse(utils.NewLogger())

	tests := []struct {
		Name                          string
		Message                       string
		ExpectedWWWAuthenticateHeader string
		ExpectedStatus                int
	}{
		{Name: "Unauthorized", Message: consts.ErrUnauthorized, ExpectedWWWAuthenticateHeader: "Bearer", ExpectedStatus: http.StatusUnauthorized},
		{Name: "Bad Request", Message: consts.ErrInvalidInput, ExpectedWWWAuthenticateHeader: "", ExpectedStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			responseHelper.ClientError(w, tt.Message, tt.ExpectedStatus)
			resp := w.Result()
			responseMessage := fmt.Sprintf(`{"message":"%s"}`, tt.Message)
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.ExpectedWWWAuthenticateHeader, resp.Header.Get("WWW-Authenticate"))
			assert.Equal(t, responseMessage, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestFailedValidation(t *testing.T) {
	t.Parallel()
	responseHelper := NewResponse(utils.NewLogger())

	tests := []struct {
		Name           string
		Errors         map[string]string
		ExpectedStatus int
	}{
		{Name: "Single error", Errors: map[string]string{"name": "invalid name"}, ExpectedStatus: http.StatusUnprocessableEntity},
		{Name: "Multiple errors", Errors: map[string]string{"name": "invalid name", "price": "invalid price"}, ExpectedStatus: http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			responseHelper.FailedValidation(w, tt.Errors)
			resp := w.Result()
			assert.Equal(t, tt.ExpectedStatus, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			// unmarshal the response body to a map to compare the errors
			var response map[string]any
			err := json.NewDecoder(resp.Body).Decode(&response)
			assert.Nil(t, err)
			assert.Equal(t, "one or more validation errors occurred", response["message"])
			assert.Equal(t, len(tt.Errors), len(response["errors"].(map[string]any)))
		})
	}
}
