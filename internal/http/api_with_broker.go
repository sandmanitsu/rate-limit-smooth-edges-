package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rate/internal/kafka"
	metric "rate/internal/metrics"
	"strconv"
	"time"
)

type HandlerProducer struct {
	producer *kafka.Producer
}

func HandlerWithQueue(producer *kafka.Producer) *http.ServeMux {
	handlerProducer := HandlerProducer{
		producer: producer,
	}

	router := http.NewServeMux()
	router.HandleFunc("POST /create", handlerProducer.create)

	return router
}

func (h *HandlerProducer) create(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if r.ContentLength == 0 {
		resp, _ := json.Marshal(Response{Result: "empty body"})
		metric.ObserveCodeStatus(400, time.Since(start))
		w.WriteHeader(400)
		_, _ = w.Write(resp)

		return
	}

	decoder := json.NewDecoder(r.Body)
	var data Request
	err := decoder.Decode(&data)
	if err != nil {
		log.Panicln(err)

		resp, _ := json.Marshal(Response{Result: "json unmarshal err"})
		metric.ObserveCodeStatus(500, time.Since(start))
		w.WriteHeader(500)
		_, _ = w.Write(resp)

		return
	}

	id, _ := strconv.Atoi(data.ProductId)
	cnt, _ := strconv.Atoi(data.Count)

	if !validate(id) || !validate(cnt) {
		resp, _ := json.Marshal(Response{Result: "invalid"})
		metric.ObserveCodeStatus(400, time.Since(start))
		w.WriteHeader(400)
		_, _ = w.Write(resp)

		return
	}

	message, _ := json.Marshal(data)
	err = h.producer.WriteMesage(context.Background(), message)
	if err != nil {
		resp, _ := json.Marshal(Response{Result: "failed produce message"})

		metric.ObserveCodeStatus(500, time.Since(start))
		w.WriteHeader(500)
		_, _ = w.Write(resp)
	}

	resp, _ := json.Marshal(Response{Result: "created"})
	metric.ObserveCodeStatus(200, time.Since(start))
	w.WriteHeader(200)
	_, _ = w.Write(resp)
}
