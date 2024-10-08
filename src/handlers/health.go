package handlers

import (
	"encoding/json"
	"github.com/sicet7/go-compose-it/src/http/middleware"
	"net/http"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (*HealthHandler) Pattern() string {
	return "/api/health"
}

func (*HealthHandler) Middleware() middleware.Middleware {
	return middleware.NewStack()
}

func (*HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"status": "Ok",
	})
	if err != nil {
		panic(err)
	}
}
