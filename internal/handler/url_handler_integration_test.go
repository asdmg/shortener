package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shortener/internal/config"
	"shortener/internal/database"
	"shortener/internal/repository"
	"shortener/internal/service"
	"testing"
)

func TestURLHandler_CreateIntegration(t *testing.T) {

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

	urlRepository := repository.NewURLRepository(db)

	urlService := service.NewURLService(
		urlRepository,
	)

	urlHandler := NewURLHandler(urlService, cfg.App.BaseURL)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /api/urls",
		urlHandler.Create,
	)

	server := httptest.NewServer(mux)

	defer server.Close()

	requestBody := map[string]string{
		"url": "https://google.com",
	}

	body, err := json.Marshal(requestBody)

	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		server.URL+"/api/urls",
		bytes.NewReader(body),
	)

	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {

		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			response.StatusCode,
		)
	}
}
