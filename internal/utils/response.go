package utils

import (
	"encoding/json"
	"net/http"
)



type GenericResponse struct {
	Code int `json:"code"`
	Status string `json:"status"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}


func NewJSONResponse(w http.ResponseWriter, code int, status string, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(GenericResponse{
		Code: code,
		Status: status,
		Message: message,
		Data: data,
	})
}