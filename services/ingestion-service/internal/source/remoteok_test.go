package source

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchJobs_Success(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`[
            {"legal": "some terms"},
            {
				"id": "123",
				"company": "Acme",
				"position": "Backend Engineer",
				"url": "https://acme.com/1",
				"location": "Remote",
				"salary_min": 50000,
				"salary_max": 90000,
				"date": "2026-01-01T00:00:00Z",
				"description": "Go and PostgreSQL"
			}
        ]`))
    }))
    defer server.Close()

    source := NewRemoteOKSource(slog.New(slog.DiscardHandler), server.Client(), server.URL)

    jobs, err := source.FetchJobs(context.Background())

    require.NoError(t, err)
    require.Len(t, jobs, 1)
    assert.Equal(t, "Acme", jobs[0].Company)
    assert.Equal(t, "Backend Engineer", jobs[0].Role)
	assert.Equal(t, "https://acme.com/1", jobs[0].OfferUrl)
	assert.Equal(t, "Remote", jobs[0].Location)
	assert.Equal(t, int32(50000), *jobs[0].SalaryMin)
	assert.Equal(t, int32(90000), *jobs[0].SalaryMax)
	assert.Equal(t, "Go and PostgreSQL", jobs[0].Description)
}

func TestFetchJobs_InvalidId(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`[
            {"legal": "some terms"},
            {
				"id": "",
				"company": "Acme",
				"position": "Backend Engineer",
				"url": "https://acme.com/1",
				"location": "Remote",
				"salary_min": 50000,
				"salary_max": 90000,
				"date": "2026-01-01T00:00:00Z",
				"description": "Go and PostgreSQL"
			}
        ]`))
    }))
    defer server.Close()

    source := NewRemoteOKSource(slog.New(slog.DiscardHandler), server.Client(), server.URL)

    jobs, err := source.FetchJobs(context.Background())

    require.NoError(t, err)
    require.Len(t, jobs, 0)
}

func TestFetchJobs_HttpError(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer server.Close()

    source := NewRemoteOKSource(slog.New(slog.DiscardHandler), server.Client(), server.URL)
    _, err := source.FetchJobs(context.Background())

    assert.Error(t, err)
}
