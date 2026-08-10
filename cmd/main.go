package main

import (
	"log"
	"net/http"

	"shortener/internal/database"
	"shortener/internal/handler"
	"shortener/internal/repository"
	"shortener/internal/service"
)

func main() {

	db, err := database.NewPostgres()

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
		"GET /{code}",
		urlHandler.Redirect,
	)

	log.Println("Server running on :8080")

	if err := http.ListenAndServe(
		":8080",
		mux,
	); err != nil {
		log.Fatal(err)
	}
}
