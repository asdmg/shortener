package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shortener/internal/model"
	"strings"
	"testing"
)

type mockURLService struct {
	createFunc func(
		ctx context.Context,
		originalURL string,
	) (*model.URL, error)

	findByCodeFunc func(
		ctx context.Context,
		code string,
	) (*model.URL, error)
}

func (m *mockURLService) Create(
	ctx context.Context,
	originalURL string,
) (*model.URL, error) {

	return m.createFunc(ctx, originalURL)
}

func (m *mockURLService) FindByCode(
	ctx context.Context,
	code string,
) (*model.URL, error) {

	return m.findByCodeFunc(ctx, code)
}

func TestURLHandler_Create(t *testing.T) {

	mockService := &mockURLService{
		createFunc: func(
			ctx context.Context,
			originalURL string,
		) (*model.URL, error) {

			return &model.URL{
				ID:          1,
				Code:        "abc123",
				OriginalURL: originalURL,
			}, nil
		},
	}

	handler := NewURLHandler(mockService)

	req := httptest.NewRequest(
		"POST",
		"/api/urls",
		strings.NewReader(
			`{"url":"https://google.com"}`,
		),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.Create(
		recorder,
		req,
	)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}
}

func TestURLHandler_Create_InvalidJSON(t *testing.T) {

	mockService := &mockURLService{
		createFunc: func(
			ctx context.Context,
			originalURL string,
		) (*model.URL, error) {

			t.Fatal("Create should not be called")

			return nil, nil
		},
	}

	handler := NewURLHandler(mockService)

	req := httptest.NewRequest(
		"POST",
		"/api/urls",
		strings.NewReader(
			`{"url":`,
		),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.Create(
		recorder,
		req,
	)

	var response ErrorResponse

	if err := json.NewDecoder(
		recorder.Body,
	).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Error != "invalid_json" {
		t.Fatalf(
			"expected error invalid_json, got %s",
			response.Error,
		)
	}

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}
