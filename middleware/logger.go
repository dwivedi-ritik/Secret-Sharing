package middleware

import (
	"log/slog"
	"net/http"
)

func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		slog.Info(r.Method, "Path", r.URL.Path)
	})
}
