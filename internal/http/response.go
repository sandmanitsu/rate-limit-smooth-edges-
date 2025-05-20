package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func response(w http.ResponseWriter, response Response) {
	resp, _ := json.Marshal(response)
	w.WriteHeader(response.Code)
	_, _ = w.Write(resp)
}
