package repository

import (
	"context"
	"shortener/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	db *pgxpool.Pool
}

func NewURLRepository(db *pgxpool.Pool) *URLRepository {
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
			original_url,
		    expires_at
		)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		url.Code,
		url.OriginalURL,
		url.ExpiresAt,
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

	err := r.db.QueryRow(
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

func (r *URLRepository) IncrementClicks(
	ctx context.Context,
	code string,
) error {

	query := `
		UPDATE urls
		SET clicks = clicks + 1
		WHERE code = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		code,
	)

	return err
}
