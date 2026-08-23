package application

import (
	"encoding/json"
	"jobradar/tracker-service/internal/domain"
	"jobradar/tracker-service/internal/handler/responses"
	"jobradar/tracker-service/internal/repository"
	"log/slog"
	"net/http"
)

type createApplicationRequest struct {
    Company string `json:"company"`
    Role    string `json:"role"`
    URL     string `json:"url"`
}

type ApplicationHandler struct {
	logger *slog.Logger
	repository *repository.ApplicationRepository
}

func NewApplicationHandler(logger *slog.Logger, repo *repository.ApplicationRepository) *ApplicationHandler {
	return &ApplicationHandler {
		logger: logger,
		repository: repo,
	}
}

func (a* ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	applicationReq := createApplicationRequest{}
	err := json.NewDecoder(r.Body).Decode(&applicationReq)

	if err != nil {
		a.logger.Warn(responses.BadRequestHandler(w, r))
		return
	}

	application, err := domain.NewApplication(applicationReq.Company, applicationReq.Role, applicationReq.URL)

	if err != nil {
		a.logger.Warn(responses.UnprocessableEntityHandler(w, r))
		return
	}

	err = a.repository.Create(r.Context(), application)

	if err != nil {
		a.logger.Error(responses.InternalServerErrorHandler(w, r))
		return
	}

	jsonRecord, err := json.Marshal(application)

	if err != nil {
		a.logger.Error(responses.InternalServerErrorHandler(w, r))
	} else {
		a.logger.Info(responses.CreatedHandler(w, r, "application/json", jsonRecord))
	}
}

func (a* ApplicationHandler) ListApplications(w http.ResponseWriter, r *http.Request) {
	lstApplications, err := a.repository.List(r.Context())

	if err != nil {
		a.logger.Error(responses.InternalServerErrorHandler(w, r))
		return
	}

	jsonRecords, err := json.Marshal(lstApplications)

	if err != nil {
		a.logger.Error(responses.InternalServerErrorHandler(w, r))
	} else {
		a.logger.Info(responses.OkHandler(w, r, "application/json", jsonRecords))
	}
}

func (h *ApplicationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// mux.HandleFunc("POST /applications", ...)
}
