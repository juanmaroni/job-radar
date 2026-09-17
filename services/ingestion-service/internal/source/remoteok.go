package source

import (
	"context"
	"encoding/json"
	"fmt"
	"jobradar/ingestion-service/internal/domain"
	"log/slog"
	"net/http"
	"time"
)

const RemoteOkApiUrl = "https://remoteok.com/api"

type remoteOKJob struct {
    ID        string `json:"id"`
    Company   string `json:"company"`
    Position  string `json:"position"`
    URL       string `json:"url"`
    Location  string `json:"location"`
    SalaryMin int32  `json:"salary_min"`
    SalaryMax int32  `json:"salary_max"`
	Description string `json:"description"`
    Date      string `json:"date"`
}

type RemoteOKSource struct {
	logger *slog.Logger
	cli *http.Client
	apiUrl string
}

func NewRemoteOKSource(logger*slog.Logger, cli *http.Client, apiUrl string) *RemoteOKSource {
    return &RemoteOKSource{
		logger: logger,
		cli: cli,
		apiUrl: apiUrl,
	}
}

func (roks *RemoteOKSource) FetchJobs(ctx context.Context) ([]domain.Job, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, roks.apiUrl, nil)

	if err != nil {
		roks.logger.Error(err.Error())

		return nil, err
	}

	resp, err := roks.cli.Do(req)

	if err != nil {
		roks.logger.Error(err.Error())

		return nil, err
	}

	if resp.StatusCode != 200 {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		roks.logger.Error(err.Error())
		
		return nil, err
	}

	defer resp.Body.Close()

	var raw []remoteOKJob
	
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		roks.logger.Error(err.Error())

		return nil, err
	}

	if len(raw) < 2 {
		roks.logger.Warn("no json data received")

		return nil, nil
	}

	var jobs []domain.Job

	for _, jsonJob := range raw[1:] {
		// Check ID and company
		if jsonJob.ID == "" || jsonJob.Company == "" {
			continue
		}

		// Parse date
		parsedDate, err := time.Parse(time.RFC3339, jsonJob.Date)
		var publishedAt *time.Time
		
		if err == nil {
			publishedAt = &parsedDate
		} else {
			roks.logger.Warn("failed to parse published date", slog.String("date", jsonJob.Date), slog.String("error", err.Error()))
			publishedAt = nil
		}

		// Normalize salaries
		salaryMin := &jsonJob.SalaryMin
		salaryMax := &jsonJob.SalaryMax

		if jsonJob.SalaryMin == 0 {
			salaryMin = nil
		}

		if jsonJob.SalaryMax == 0 {
			salaryMax = nil
		}

		// Construct job
		job, err := domain.NewJob("RemoteOK", jsonJob.ID, jsonJob.Company, jsonJob.Position, jsonJob.URL, jsonJob.Location, jsonJob.Description, salaryMin, salaryMax, publishedAt)

		if err != nil {
			roks.logger.Warn("failed to construct job struct", slog.String("error", err.Error()))
			continue
		}

		jobs = append(jobs, *job)
	}

	return jobs, nil
}
