package main

import (
	"log"
	"net/http"
	api "rate/internal/http"
	metric "rate/internal/metrics"
)

const (
	serverHost = "localhost:8083"
	metricHost = "localhost:8084"
)

func main() {
	go func() {
		err := metric.Listen(metricHost)
		if err != nil {
			panic(err)
		}
	}()

	handler := api.Handler()

	s := &http.Server{
		Addr:    serverHost,
		Handler: handler,
	}

	log.Fatal(s.ListenAndServe())
}
