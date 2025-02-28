package middleware

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/daut/jed/cmd/api/consts"
)

func (mw *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			mw.Response.ClientError(w, "missing authorization header", http.StatusBadRequest)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			mw.Response.ClientError(w, consts.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		hash := sha256.Sum256([]byte(headerParts[1]))

		token, err := mw.Queries.GetToken(r.Context(), hash[:])
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				mw.Response.ClientError(w, consts.ErrUnauthorized, http.StatusUnauthorized)
			} else {
				mw.Response.ServerError(w, err)
			}
			return
		}

		if token.ExpiresAt.Time.Before(time.Now()) {
			mw.Response.ClientError(w, "expired token", http.StatusUnauthorized)
			return
		}

		admin, err := mw.Queries.GetAdminByID(r.Context(), token.AdminID)
		if err != nil {
			mw.Response.ServerError(w, err)
		}

		r = r.WithContext(context.WithValue(r.Context(), "admin", admin))

		next.ServeHTTP(w, r)
	})
}
