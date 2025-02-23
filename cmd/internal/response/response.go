package response

import (
	"encoding/json"
	"net/http"
)

const (
	StatusSuccess = "success"
	StatusError = "error"
)

type Response struct {
	Status string `json:"status"`
	Message string `json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, val interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(val)
}

func Success(w http.ResponseWriter, statusCode int, message string, val interface{}) error {
	return WriteJSON(w, statusCode, Response{
		Status: StatusSuccess,
		Message: message,
		Data: val,
	})
}

func Error(w http.ResponseWriter, statusCode int, message string, errors ...string) error {
	return WriteJSON(w, statusCode, Response{
		Status: StatusError,
		Message: message,
		Errors: errors,
	})
}
