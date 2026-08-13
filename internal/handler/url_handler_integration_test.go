package handler

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shortener/internal/config"
	"shortener/internal/database"
	"shortener/internal/model"
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

func TestURLHandler_RedirectIntegration(t *testing.T) {

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

	urlHandler := NewURLHandler(
		urlService,
		cfg.App.BaseURL,
	)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /{code}",
		urlHandler.Redirect,
	)

	server := httptest.NewServer(mux)

	defer server.Close()

	ctx := context.Background()

	testURL := &model.URL{
		Code:        testCode(t),
		OriginalURL: "https://google.com",
	}

	err = urlRepository.Create(
		ctx,
		testURL,
	)

	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest(
		http.MethodGet,
		server.URL+"/"+testURL.Code,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{
		CheckRedirect: func(
			req *http.Request,
			via []*http.Request,
		) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := client.Do(request)

	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusFound {

		t.Fatalf(
			"expected status %d, got %d",
			http.StatusFound,
			response.StatusCode,
		)
	}

	location := response.Header.Get("Location")

	if location != testURL.OriginalURL {

		t.Fatalf(
			"expected Location %s, got %s",
			testURL.OriginalURL,
			location,
		)
	}

	foundURL, err := urlRepository.FindByCode(
		ctx,
		testURL.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if foundURL.Clicks != 1 {

		t.Fatalf(
			"expected 1 click, got %d",
			foundURL.Clicks,
		)
	}
}

func testCode(t *testing.T) string {
	t.Helper()

	bytes := make([]byte, 4)

	if _, err := rand.Read(bytes); err != nil {
		t.Fatal(err)
	}

	return hex.EncodeToString(bytes)[:6]
}
