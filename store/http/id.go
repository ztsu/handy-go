package http

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"
)

const QueryStringIDKey = "ID"

func QueryStringID(paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), QueryStringIDKey, chi.URLParam(r, paramName))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if id, ok := ctx.Value(QueryStringIDKey).(string); ok {
		return id
	}

	return ""
}
