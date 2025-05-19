package main

import (
	"log"
	"net/http"
	api "rate/internal/http"
)

func main() {
	hadler := api.Handler()

	s := &http.Server{
		Addr:    "localhost:8083",
		Handler: hadler,
	}

	log.Fatal(s.ListenAndServe())
}
