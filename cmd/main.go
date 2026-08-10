package main

import (
	"fmt"
	"log"
	"net/http"

	"shortener/internal/database"
)

func main() {

	db, err := database.NewPostgres()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "URL Shortener")
	})

	fmt.Println("Server running on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
