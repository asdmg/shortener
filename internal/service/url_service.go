package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/url"
	"shortener/internal/model"
	"shortener/internal/repository"
	"strings"
)

var (
	ErrInvalidURL = errors.New("invalid URL")
)

type URLService struct {
	repository *repository.URLRepository
}

func NewURLService(
	repository *repository.URLRepository,
) *URLService {

	return &URLService{
		repository: repository,
	}
}

func (s *URLService) Create(
	ctx context.Context,
	originalURL string,
) (*model.URL, error) {

	if err := validateURL(originalURL); err != nil {
		return nil, err
	}

	code, err := generateCode(6)

	if err != nil {
		return nil, err
	}

	url := &model.URL{
		Code:        code,
		OriginalURL: strings.TrimSpace(originalURL),
	}

	if err := s.repository.Create(ctx, url); err != nil {
		return nil, err
	}

	return url, nil
}

func generateCode(size int) (string, error) {

	bytes := make([]byte, size)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes)[:size], nil
}

func (s *URLService) FindByCode(
	ctx context.Context,
	code string,
) (*model.URL, error) {

	return s.repository.FindByCode(
		ctx,
		code,
	)
}

func validateURL(value string) error {
	value = strings.TrimSpace(value)

	if value == "" {
		return ErrInvalidURL
	}

	parsedURL, err := url.ParseRequestURI(value)

	if err != nil {
		return ErrInvalidURL
	}

	if parsedURL.Scheme != "http" &&
		parsedURL.Scheme != "https" {
		return ErrInvalidURL
	}

	if parsedURL.Host == "" {
		return ErrInvalidURL
	}

	return nil
}
