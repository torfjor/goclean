package authorization

import (
	"context"
	"goclean"
	"net/http"
)

func HTTPAuthorizer(next http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "password" {
			http.Error(rw, "Wrong password", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), goclean.ContextKeyPermissions, goclean.Greet)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	}
}
