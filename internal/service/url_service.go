package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"net/url"
	"shortener/internal/model"
	"shortener/internal/repository"
	"strings"
	"time"
)

var (
	ErrInvalidURL           = errors.New("invalid URL")
	ErrURLNotFound          = errors.New("url not found")
	ErrURLExpired           = errors.New("url expired")
	ErrInvalidExpiration    = errors.New("invalid expiration date")
	ErrCodeGenerationFailed = errors.New(
		"failed to generate unique URL code",
	)
)

type URLRepository interface {
	Create(
		ctx context.Context,
		url *model.URL,
	) error

	FindByCode(
		ctx context.Context,
		code string,
	) (*model.URL, error)

	IncrementClicks(
		ctx context.Context,
		code string,
	) error
}

type URLService struct {
	repository URLRepository
}

func NewURLService(
	repository URLRepository,
) *URLService {
	return &URLService{
		repository: repository,
	}
}

func (s *URLService) Create(
	ctx context.Context,
	originalURL string,
	expiresAt *time.Time,
) (*model.URL, error) {

	if err := validateURL(originalURL); err != nil {
		return nil, err
	}

	const maxCreateAttempts = 5

	for attempt := 0; attempt < maxCreateAttempts; attempt++ {

		code, err := generateCode(6)

		if err != nil {
			return nil, err
		}

		url := &model.URL{
			Code:        code,
			OriginalURL: strings.TrimSpace(originalURL),
			ExpiresAt:   expiresAt,
		}

		err = s.repository.Create(
			ctx,
			url,
		)

		if err == nil {
			return url, nil
		}

		if !errors.Is(
			err,
			repository.ErrCodeAlreadyExists,
		) {
			return nil, err
		}
	}

	return nil, ErrCodeGenerationFailed
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

	url, err := s.repository.FindByCode(
		ctx,
		code,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrURLNotFound
		}

		return nil, err
	}

	if url.ExpiresAt != nil &&
		time.Now().After(*url.ExpiresAt) {

		return nil, ErrURLExpired
	}

	return url, nil
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

func (s *URLService) IncrementClicks(
	ctx context.Context,
	code string,
) error {

	return s.repository.IncrementClicks(
		ctx,
		code,
	)
}
