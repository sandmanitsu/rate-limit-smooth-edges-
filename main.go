package main

import (
	"log"
	"net/http"
	api "rate/internal/http"
	"rate/internal/kafka"
	metric "rate/internal/metrics"
)

const (
	serverHost = "localhost:8083"
	metricHost = "localhost:8084"

	kafkaBroker = "localhost:9094"
	kafkaTopic  = "msgs"
)

func main() {
	go func() {
		err := metric.Listen(metricHost)
		if err != nil {
			panic(err)
		}
	}()

	// запускаем обработку запрос в очередь
	// записываем все запросы в очередь, потом вычитываем их используя worker pool
	producer := kafka.NewProducer(kafkaBroker, kafkaTopic)
	kafka.StartConsumer(kafkaBroker, kafkaTopic)
	handler := api.HandlerWithQueue(producer)

	// запускает обработчик в rate limit'ом в 50 rps
	// handler := api.Handler()

	s := &http.Server{
		Addr:    serverHost,
		Handler: handler,
	}

	log.Fatal(s.ListenAndServe())
}
