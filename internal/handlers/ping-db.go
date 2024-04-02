package handlers

import (
	"context"
	"net/http"
	"time"

	db "github.com/fishus/go-advanced-metrics/internal/database"
)

// PingDBHandler processes the request GET /ping.
// Checks the connection to the database.
func PingDBHandler(w http.ResponseWriter, r *http.Request) {
	dbPool, err := db.Pool()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), (3 * time.Second))
	defer cancel()

	if err := dbPool.Ping(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
