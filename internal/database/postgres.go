package database

import (
	"database/sql"
	"fmt"

	"shortener/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgres(
	cfg config.DatabaseConfig,
) (*sql.DB, error) {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
