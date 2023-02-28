package gateway

import (
	"encoding/json"
	"net/http"
)

const (
	// OK success
	OK = 0
	// RequestErr request error
	RequestErr = -400
	// ServerErr server error
	ServerErr = -500
)

type resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func errors(w http.ResponseWriter, code int, msg string) {
	r := resp{Code: code, Message: msg}
	b, _ := json.Marshal(r)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func result(w http.ResponseWriter, data interface{}, code int) {
	r := resp{Code: code, Data: data}
	b, _ := json.Marshal(r)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}
