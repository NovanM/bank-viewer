package http

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	response := APIResponse{
		Status:  false,
		Message: message,
		Data:    nil,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}

func RespondWithJSON(w http.ResponseWriter, code int, message string, data interface{}) {
	response := APIResponse{
		Status:  true,
		Message: message,
		Data:    data,
	}

	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}
