package job

import (
	"encoding/json"
	"jobradar/jobs-service/internal/domain"
	"jobradar/jobs-service/internal/handler/responses"
	"log/slog"
	"net/http"
	"strconv"
)

type JobHandler struct {
    logger *slog.Logger
    store  domain.JobStore
}

func NewJobHandler(logger *slog.Logger, store domain.JobStore) *JobHandler {
	return &JobHandler{
		logger: logger,
		store: store,
	}
}

func (jh *JobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	var days *int
	daysQuery := r.URL.Query().Get("days")

	if daysQuery != "" {
		daysParsed, err := strconv.Atoi(daysQuery)

		if err != nil || daysParsed <= 0 {
			days = nil
		} else {
			days = &daysParsed
		}
	}

	filters := domain.JobFilters{Days: days}

	lstJobs, err := jh.store.List(r.Context(), filters)

	if err != nil {
		jh.logger.Error(responses.InternalServerErrorHandler(w, r))

		return
	}

	jsonRecords, err := json.Marshal(lstJobs)

	if err != nil {
		jh.logger.Error(responses.InternalServerErrorHandler(w, r))
	} else {
		jh.logger.Info(responses.OkHandler(w, r, "application/json", jsonRecords))
	}
}
