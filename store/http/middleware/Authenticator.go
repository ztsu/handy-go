package middleware

import (
	"context"
	"net/http"
)

const UserIDKey = "userID"

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("User-ID")
		if id == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) string {
	var userID string

	if ctx == nil {
		return userID
	}

	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}

	return userID
}
