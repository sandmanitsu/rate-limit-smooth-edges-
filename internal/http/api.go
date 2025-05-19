package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	limiter := rate.NewLimiter(rateLimit, burst)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			fmt.Println(r.Method, time.Now(), " - failed")

			resp, _ := json.Marshal(Response{Result: "server kaput"})
			w.WriteHeader(500)
			w.Write(resp)

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

type Response struct {
	Result string `json:"result"`
}

func create(w http.ResponseWriter, r *http.Request) {
	// time.Sleep(500 * time.Millisecond)

	decoder := json.NewDecoder(r.Body)
	var data Request
	err := decoder.Decode(&data)
	if err != nil {
		log.Panicln(err)

		resp, _ := json.Marshal(Response{Result: "json unmarshal err"})
		w.WriteHeader(500)
		w.Write(resp)

		return
	}

	id, _ := strconv.Atoi(data.ProductId)
	cnt, _ := strconv.Atoi(data.Count)

	if !validate(id) || !validate(cnt) {
		resp, _ := json.Marshal(Response{Result: "invalid"})
		w.WriteHeader(400)
		w.Write(resp)

		return
	}

	resp, _ := json.Marshal(Response{Result: "created"})
	w.WriteHeader(200)
	w.Write(resp)
}

func validate(n int) bool {
	if n < 0 || n > 100 {
		return false
	}

	return true
}
