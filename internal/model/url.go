package model

import "time"

type URL struct {
	ID          int64
	Code        string
	OriginalURL string
	Clicks      int64
	CreatedAt   time.Time
	ExpiresAt   *time.Time
}
