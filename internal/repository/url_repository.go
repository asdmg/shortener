package repository

import (
	"context"
	"database/sql"

	"shortener/internal/model"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) Create(
	ctx context.Context,
	url *model.URL,
) error {

	query := `
		INSERT INTO urls (
			code,
			original_url
		)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		url.Code,
		url.OriginalURL,
	).Scan(
		&url.ID,
		&url.CreatedAt,
	)
}

func (r *URLRepository) FindByCode(
	ctx context.Context,
	code string,
) (*model.URL, error) {

	query := `
		SELECT
			id,
			code,
			original_url,
			clicks,
			created_at,
			expires_at
		FROM urls
		WHERE code = $1
	`

	url := &model.URL{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		code,
	).Scan(
		&url.ID,
		&url.Code,
		&url.OriginalURL,
		&url.Clicks,
		&url.CreatedAt,
		&url.ExpiresAt,
	)

	if err != nil {
		return nil, err
	}

	return url, nil
}
