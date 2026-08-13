package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shortener/internal/config"
	"shortener/internal/database"
	"shortener/internal/handler"
	"shortener/internal/repository"
	"shortener/internal/service"
	"syscall"
	"time"
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
		cfg.App.BaseURL,
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

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	serverErr := make(chan error, 1)

	go func() {

		log.Println("Server running on", addr)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			serverErr <- err
		}

	}()

	select {

	case err := <-serverErr:

		log.Fatalf(
			"server error: %v",
			err,
		)

	case <-stop:

		log.Println("Shutting down server...")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {

		log.Printf(
			"server shutdown error: %v",
			err,
		)
	}

	log.Println("Server stopped")
}
