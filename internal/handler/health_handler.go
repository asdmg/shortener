package handler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) Health(
	w http.ResponseWriter,
	r *http.Request,
) {
	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "ok",
		},
	)
}

func (h *HealthHandler) Ready(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := h.db.Ping(r.Context()); err != nil {
		writeError(
			w,
			http.StatusServiceUnavailable,
			"database_unavailable",
			"database is unavailable",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "ready",
		},
	)
}
