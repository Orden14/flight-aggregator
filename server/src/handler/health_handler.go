package handler

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (healthHandler *HealthHandler) ServeHTTP(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")

	writer.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(writer).Encode(map[string]string{"status": "ok"})
}
