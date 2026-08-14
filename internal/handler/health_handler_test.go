package handler

import (
	"net/http"
	"net/http/httptest"
	"shortener/internal/config"
	"shortener/internal/database"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestHealthHandler_Health(t *testing.T) {

	handler := NewHealthHandler(
		&pgxpool.Pool{},
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	response := httptest.NewRecorder()

	handler.Health(
		response,
		request,
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}
}

func TestHealthHandler_Ready(t *testing.T) {

	cfg, err := config.Load()

	if err != nil {
		t.Fatal(err)
	}

	db, err := database.NewPostgres(
		cfg.Database,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	handler := NewHealthHandler(db)

	request := httptest.NewRequest(
		http.MethodGet,
		"/ready",
		nil,
	)

	response := httptest.NewRecorder()

	handler.Ready(
		response,
		request,
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}
}
