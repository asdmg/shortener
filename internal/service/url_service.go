package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"shortener/internal/model"
	"shortener/internal/repository"
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

	code, err := generateCode(6)

	if err != nil {
		return nil, err
	}

	url := &model.URL{
		Code:        code,
		OriginalURL: originalURL,
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
