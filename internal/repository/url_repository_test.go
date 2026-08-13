package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"shortener/internal/config"
	"shortener/internal/database"
	"shortener/internal/model"
	"sync"
	"testing"
)

func TestURLRepository_Create(t *testing.T) {

	cfg, err := config.Load()

	if err != nil {
		t.Fatal(err)
	}

	db, err := database.NewPostgres(
		cfg.Database,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	repository := NewURLRepository(db)

	url := &model.URL{
		Code:        testCode(t),
		OriginalURL: "https://google.com",
	}

	err = repository.Create(
		context.Background(),
		url,
	)

	if err != nil {
		t.Fatal(err)
	}

	if url.ID == 0 {
		t.Fatal("expected URL ID to be generated")
	}

	if url.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be generated")
	}
}

func TestURLRepository_FindByCode(t *testing.T) {

	cfg, err := config.Load()

	if err != nil {
		t.Fatal(err)
	}

	db, err := database.NewPostgres(
		cfg.Database,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	repository := NewURLRepository(db)

	url := &model.URL{
		Code:        testCode(t),
		OriginalURL: "https://google.com",
	}

	err = repository.Create(
		context.Background(),
		url,
	)

	if err != nil {
		t.Fatal(err)
	}

	foundURL, err := repository.FindByCode(
		context.Background(),
		url.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if foundURL.Code != url.Code {
		t.Fatalf(
			"expected code %s, got %s",
			url.Code,
			foundURL.Code,
		)
	}

	if foundURL.OriginalURL != url.OriginalURL {
		t.Fatalf(
			"expected URL %s, got %s",
			url.OriginalURL,
			foundURL.OriginalURL,
		)
	}
}

func TestURLRepository_IncrementClicks(t *testing.T) {

	cfg, err := config.Load()

	if err != nil {
		t.Fatal(err)
	}

	db, err := database.NewPostgres(
		cfg.Database,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	repository := NewURLRepository(db)

	ctx := context.Background()

	url := &model.URL{
		Code:        testCode(t),
		OriginalURL: "https://google.com",
	}

	err = repository.Create(ctx, url)

	if err != nil {
		t.Fatal(err)
	}

	if url.Clicks != 0 {
		t.Fatalf(
			"expected clicks to be 0, got %d",
			url.Clicks,
		)
	}

	err = repository.IncrementClicks(
		ctx,
		url.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	foundURL, err := repository.FindByCode(
		ctx,
		url.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if foundURL.Clicks != 1 {
		t.Fatalf(
			"expected clicks to be 1, got %d",
			foundURL.Clicks,
		)
	}
}

func TestURLRepository_IncrementClicksConcurrent(t *testing.T) {

	cfg, err := config.Load()

	if err != nil {
		t.Fatal(err)
	}

	db, err := database.NewPostgres(
		cfg.Database,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	repository := NewURLRepository(db)

	ctx := context.Background()

	url := &model.URL{
		Code:        testCode(t),
		OriginalURL: "https://google.com",
	}

	if err := repository.Create(ctx, url); err != nil {
		t.Fatal(err)
	}

	const numberOfClicks = 100

	var wg sync.WaitGroup

	wg.Add(numberOfClicks)

	for i := 0; i < numberOfClicks; i++ {

		go func() {

			defer wg.Done()

			if err := repository.IncrementClicks(
				ctx,
				url.Code,
			); err != nil {

				t.Errorf(
					"failed to increment clicks: %v",
					err,
				)
			}

		}()
	}

	wg.Wait()

	foundURL, err := repository.FindByCode(
		ctx,
		url.Code,
	)

	if err != nil {
		t.Fatal(err)
	}

	if foundURL.Clicks != numberOfClicks {

		t.Fatalf(
			"expected %d clicks, got %d",
			numberOfClicks,
			foundURL.Clicks,
		)
	}
}

func testCode(t *testing.T) string {

	t.Helper()

	bytes := make([]byte, 4)

	if _, err := rand.Read(bytes); err != nil {
		t.Fatal(err)
	}

	return hex.EncodeToString(bytes)[:6]
}
