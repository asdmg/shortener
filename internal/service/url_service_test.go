package service

import (
	"context"
	"errors"
	"shortener/internal/model"
	"shortener/internal/repository"
	"testing"
)

type fakeURLRepository struct {
	createCalls     int
	failFirstCreate bool
}

func (f *fakeURLRepository) Create(
	ctx context.Context,
	url *model.URL,
) error {

	f.createCalls++

	if f.failFirstCreate && f.createCalls == 1 {
		return repository.ErrCodeAlreadyExists
	}

	return nil
}

func (f *fakeURLRepository) FindByCode(
	ctx context.Context,
	code string,
) (*model.URL, error) {

	return nil, errors.New("not implemented")
}

func (f *fakeURLRepository) IncrementClicks(
	ctx context.Context,
	code string,
) error {

	return nil
}

func TestURLService_Create(t *testing.T) {

	repository := &fakeURLRepository{}

	service := NewURLService(repository)

	url, err := service.Create(
		context.Background(),
		"https://google.com",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	if url == nil {
		t.Fatal("expected URL")
	}

	if url.Code == "" {
		t.Fatal("expected generated code")
	}

	if repository.createCalls != 1 {
		t.Fatalf(
			"expected 1 repository call, got %d",
			repository.createCalls,
		)
	}
}

func TestURLService_Create_RetryOnCodeCollision(t *testing.T) {

	repository := &fakeURLRepository{
		failFirstCreate: true,
	}

	service := NewURLService(repository)

	url, err := service.Create(
		context.Background(),
		"https://google.com",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	if url == nil {
		t.Fatal("expected URL")
	}

	if repository.createCalls != 2 {
		t.Fatalf(
			"expected 2 repository calls, got %d",
			repository.createCalls,
		)
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid https url",
			url:     "https://google.com",
			wantErr: false,
		},
		{
			name:    "valid http url",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "empty url",
			url:     "",
			wantErr: true,
		},
		{
			name:    "domain without scheme",
			url:     "google.com",
			wantErr: true,
		},
		{
			name:    "invalid scheme",
			url:     "ftp://google.com",
			wantErr: true,
		},
		{
			name:    "invalid url",
			url:     "abc",
			wantErr: true,
		},
		{
			name:    "missing host",
			url:     "https://",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := validateURL(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"validateURL() error = %v, wantErr = %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
