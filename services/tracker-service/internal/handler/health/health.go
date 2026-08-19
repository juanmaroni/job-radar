package health

import (
	"jobradar/tracker-service/internal/handler/responses"
	"log/slog"
	"net/http"
)

type HealthHandler struct {
	logger *slog.Logger
}

func NewHealthHandler(logger *slog.Logger) *HealthHandler {
    return &HealthHandler{logger: logger}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responses.OkHandler(w, r, []byte("ok"))
}
