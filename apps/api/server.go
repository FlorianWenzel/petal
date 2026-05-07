package main

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/jackc/pgx/v5/pgxpool"
)

type server struct {
	pool *pgxpool.Pool
	cfg  *config
}

func newServer(pool *pgxpool.Pool, cfg *config) http.Handler {
	s := &server{pool: pool, cfg: cfg}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("POST /auth/signup", s.signup)
	mux.HandleFunc("POST /auth/login", s.login)
	mux.HandleFunc("POST /auth/logout", s.logout)
	mux.Handle("GET /auth/me", s.authMiddleware(http.HandlerFunc(s.me)))
	mux.Handle("PATCH /auth/me/username", s.authMiddleware(http.HandlerFunc(s.rename)))
	mux.Handle("GET /entries/{year}/{month}", s.authMiddleware(http.HandlerFunc(s.loadMonth)))
	mux.Handle("PUT /entries/{year}/{month}/{day}", s.authMiddleware(http.HandlerFunc(s.saveDay)))

	return s.corsMiddleware(mux)
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && slices.Contains(s.cfg.CORSOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
		}
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cookie")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"message": msg})
}

func decodeJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
