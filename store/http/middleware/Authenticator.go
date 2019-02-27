package middleware

import (
	"context"
	"github.com/google/uuid"
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

		ctx := context.WithValue(r.Context(), "userID", uuid.UUID(id))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) uuid.UUID {
	var userID uuid.UUID

	if ctx == nil {
		return userID
	}

	if id, ok := ctx.Value(UserIDKey).(uuid.UUID); ok {
		return id
	}

	return userID
}