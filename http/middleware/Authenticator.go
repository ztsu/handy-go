package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/ztsu/handy-go"
	"net/http"
)

const UserIDKey = "userID"

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get("User-ID")
		id, err:= uuid.Parse(userIDStr)
		if err != nil {
			if userIDStr == "" {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}

		ctx := context.WithValue(r.Context(), "userID", handy.UUID(id))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) handy.UUID {
	var userID handy.UUID

	if ctx == nil {
		return userID
	}

	if uid, ok := ctx.Value(UserIDKey).(handy.UUID); ok {
		return uid
	}

	return userID
}