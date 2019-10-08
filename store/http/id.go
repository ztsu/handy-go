package http

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"
)

const idCtxKey = "ID"

func idCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), idCtxKey, chi.URLParam(r, "ID"))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getIDCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if id, ok := ctx.Value(idCtxKey).(string); ok {
		return id
	}

	return ""
}
