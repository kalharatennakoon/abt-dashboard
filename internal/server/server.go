package server

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"abt-dashboard/internal/handlers"
)

type Server struct {
	mux *http.ServeMux
}

// gzipResponseWriter wraps http.ResponseWriter to provide gzip compression
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// gzipMiddleware provides gzip compression for responses
func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if client accepts gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Set gzip headers
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Create gzip writer
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Wrap response writer
		gzipWriter := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzipWriter, r)
	})
}

func New(api *handlers.API, staticDir string) *Server {
	mux := http.NewServeMux()

	// API routes with gzip compression
	mux.Handle("GET /api/revenue/countries", gzipMiddleware(http.HandlerFunc(api.CountryRevenue)))
	mux.Handle("GET /api/products/top", gzipMiddleware(http.HandlerFunc(api.TopProducts)))
	mux.Handle("GET /api/sales/by-month", gzipMiddleware(http.HandlerFunc(api.SalesByMonth)))
	mux.Handle("GET /api/regions/top", gzipMiddleware(http.HandlerFunc(api.TopRegions)))

	// Serve static frontend if needed
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("/", fs)

	return &Server{mux: mux}
}

func (s *Server) Listen(addr string) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           s.mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return srv.ListenAndServe()
}
