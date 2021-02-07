package main

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
)

type contextKey int

const (
	dbContextKey contextKey = iota
)

func DBMiddleware(db *Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), dbContextKey, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func DBFromRequest(r *http.Request) *Database {
	db, ok := r.Context().Value(dbContextKey).(*Database)
	if !ok {
		panic("database not in request context")
	}
	return db
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func HandleStatic(fileName, contentType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fileName == "" {
			fileName = r.URL.Path
		}
		if strings.Contains(fileName, "..") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		assetPath := filepath.Join("client", "build", fileName)

		w.Header().Set("Content-Type", contentType)
		http.ServeFile(w, r, assetPath)
	})
}
