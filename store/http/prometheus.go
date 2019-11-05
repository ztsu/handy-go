package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "handy_store",
			Name:      "requests_total",
			Help:      "The total number of requests",
		},
		[]string{"method", "path", "code"},
	)

	requestTimer = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "handy_store",
			Name:      "requests_duration_seconds",
		},
		[]string{"method", "path", "code"},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestTimer)
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()

		defer func() {
			ctx := chi.RouteContext(r.Context())
			lvs := []string{
				r.Method,
				strings.TrimRight(ctx.RoutePattern(), "/"),
				strconv.Itoa(ww.Status()),
			}

			requestsTotal.WithLabelValues(lvs...).Inc()
			requestTimer.WithLabelValues(lvs...).Observe(time.Since(start).Seconds())
		}()

		next.ServeHTTP(ww, r)
	}

	return http.HandlerFunc(fn)
}
