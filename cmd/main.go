package main

import (
	"log"
	"net/http"
	"shortener/internal/config"
	"shortener/internal/service"

	"shortener/internal/database"
	"shortener/internal/handler"
	"shortener/internal/repository"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgres(cfg.Database)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	urlRepository := repository.NewURLRepository(db)

	urlService := service.NewURLService(
		urlRepository,
	)

	urlHandler := handler.NewURLHandler(
		urlService,
	)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /api/urls",
		urlHandler.Create,
	)

	mux.HandleFunc(
		"GET /api/urls/{code}",
		urlHandler.Get,
	)

	mux.HandleFunc(
		"GET /{code}",
		urlHandler.Redirect,
	)

	addr := ":" + cfg.App.Port

	log.Println("Server running on", addr)

	if err := http.ListenAndServe(
		addr,
		mux,
	); err != nil {
		log.Fatal(err)
	}
}
