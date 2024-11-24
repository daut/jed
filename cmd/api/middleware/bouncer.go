package middleware

import (
	"net/http"

	"github.com/daut/jed/sqlc"
)

func (mw *Middleware) RequireAdminUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value("admin").(sqlc.Admin)
		if !ok {
			mw.Response.ClientError(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
