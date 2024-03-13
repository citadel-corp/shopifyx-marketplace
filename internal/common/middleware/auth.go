package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/jwt"
)

type ContextAuthKey struct{}

func Authenticate(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusForbidden)
			slog.InfoContext(r.Context(), "Missing authorization header")
			return
		}

		tokenString = tokenString[len("Bearer "):]

		subject, err := jwt.VerifyAndGetSubject(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			slog.InfoContext(r.Context(), "Invalid token: %v", err)
			return
		}

		ctx := context.WithValue(r.Context(), ContextAuthKey{}, subject)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
