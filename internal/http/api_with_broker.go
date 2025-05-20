package api

import (
	"context"
	"encoding/json"
	"net/http"
	"rate/internal/kafka"
	metric "rate/internal/metrics"
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
	var resp Response
	start := time.Now()

	defer func() {
		metric.ObserveCodeStatus(400, time.Since(start))
		response(w, resp)
	}()

	if r.ContentLength == 0 {
		resp = Response{
			Code:    http.StatusBadRequest,
			Message: "empty body",
		}

		return
	}

	decoder := json.NewDecoder(r.Body)
	var data Request
	err := decoder.Decode(&data)
	if err != nil {
		resp = Response{
			Code:    http.StatusInternalServerError,
			Message: "json unmarshal err",
		}

		return
	}

	if !validate(data) {
		resp = Response{
			Code:    http.StatusBadRequest,
			Message: "invalid data",
		}

		return
	}

	message, _ := json.Marshal(data)
	err = h.producer.WriteMesage(context.Background(), message)
	if err != nil {
		resp = Response{
			Code:    http.StatusBadRequest,
			Message: "failed produce message",
		}

		return
	}

	resp = Response{
		Code:    http.StatusOK,
		Message: "created",
	}
}
