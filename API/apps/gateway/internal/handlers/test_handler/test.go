package test_handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// TestHandler обрабатывает тестовые эндпоинты
type TestHandler struct{}

// NewTestHandler создает новый экземпляр TestHandler
func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

// HealthResponse представляет ответ health check эндпоинта
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// GetHealth обрабатывает GET /health запрос
func (h *TestHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
