package application

import (
	"encoding/json"
	"errors"
	"jobradar/tracker-service/internal/domain"
	"jobradar/tracker-service/internal/handler/responses"
	"jobradar/tracker-service/internal/repository"
	"log/slog"
	"net/http"
	"strconv"
)

type createApplicationRequest struct {
    Company string `json:"company"`
    Role string `json:"role"`
    URL string `json:"url"`
	MatchScore *int32 `json:"matchScore"`
}

type updateStatusRequest struct {
    Status domain.Status `json:"status"`
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

	application, err := domain.NewApplication(applicationReq.Company, applicationReq.Role, applicationReq.URL, applicationReq.MatchScore)

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

func (ah *ApplicationHandler) UpdateStatusApplication(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	
	if err != nil {
		ah.logger.Warn(responses.BadRequestHandler(w, r))

		return
	}

	updateStatusReq := updateStatusRequest{}
	err = json.NewDecoder(r.Body).Decode(&updateStatusReq)

	if err != nil {
		ah.logger.Warn(responses.BadRequestHandler(w, r))
		
		return
	}

	if !updateStatusReq.Status.Valid() {
		ah.logger.Warn(responses.UnprocessableEntityHandler(w, r))
		
		return
	}

	err = ah.repository.UpdateStatus(r.Context(), id, updateStatusReq.Status)

	if errors.Is(err, repository.ErrNotFound) {
		ah.logger.Warn(responses.NotFoundHandler(w, r))
		
		return
	} else if err != nil {
		ah.logger.Error(responses.InternalServerErrorHandler(w, r))

		return
	}

	ah.logger.Info(responses.NoContentHandler(w, r))
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
