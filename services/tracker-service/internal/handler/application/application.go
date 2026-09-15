package application

import (
	"encoding/json"
	"jobradar/tracker-service/internal/domain"
	"jobradar/tracker-service/internal/handler/responses"
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
	repository domain.ApplicationStore
}

func NewApplicationHandler(logger *slog.Logger, repo domain.ApplicationStore) *ApplicationHandler {
	return &ApplicationHandler {
		logger: logger,
		repository: repo,
	}
}

func (ah *ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	applicationReq := createApplicationRequest{}
	err := json.NewDecoder(r.Body).Decode(&applicationReq)

	if err != nil {
		ah.logger.Warn(responses.BadRequestHandler(w, r))
		
		return
	}

	application, err := domain.NewApplication(applicationReq.Company, applicationReq.Role, applicationReq.URL)

	if err != nil {
		ah.logger.Warn(responses.UnprocessableEntityHandler(w, r))

		return
	}

	err = ah.repository.Create(r.Context(), application)

	if err != nil {
		ah.logger.Error(responses.InternalServerErrorHandler(w, r))

		return
	}

	jsonRecord, err := json.Marshal(application)

	if err != nil {
		ah.logger.Error(responses.InternalServerErrorHandler(w, r))
	} else {
		ah.logger.Info(responses.CreatedHandler(w, r, "application/json", jsonRecord))
	}
}

func (ah *ApplicationHandler) ListApplications(w http.ResponseWriter, r *http.Request) {
	lstApplications, err := ah.repository.List(r.Context())

	if err != nil {
		ah.logger.Error(responses.InternalServerErrorHandler(w, r))

		return
	}

	jsonRecords, err := json.Marshal(lstApplications)

	if err != nil {
		ah.logger.Error(responses.InternalServerErrorHandler(w, r))
	} else {
		ah.logger.Info(responses.OkHandler(w, r, "application/json", jsonRecords))
	}
}
