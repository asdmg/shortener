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
