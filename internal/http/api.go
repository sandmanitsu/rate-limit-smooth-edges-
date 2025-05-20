package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	metric "rate/internal/metrics"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

const (
	rateLimit = 50
	burst     = 50
)

func Handler() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("POST /create", rateLimiter(create))

	return router
}

func rateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	start := time.Now()

	limiter := rate.NewLimiter(rateLimit, burst)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			fmt.Println(r.Method, time.Now(), " - failed")

			metric.ObserveCodeStatus(500, time.Since(start))
			response(w, Response{
				Code:    http.StatusTooManyRequests,
				Message: "to many request",
			})

			return
		}

		fmt.Println(r.Method, time.Now(), " - allowed")

		next(w, r)
	})
}

type Request struct {
	ProductId string `json:"product_id"`
	Count     string `json:"count"`
	Username  string `json:"username"`
}

func create(w http.ResponseWriter, r *http.Request) {
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

	resp = Response{
		Code:    http.StatusOK,
		Message: "created",
	}
}

func validate(data Request) bool {
	id, _ := strconv.Atoi(data.ProductId)
	cnt, _ := strconv.Atoi(data.Count)

	if id < 0 || id > 100 {
		return false
	}

	if cnt < 0 || cnt > 100 {
		return false
	}

	return true
}
