package application

import (
	"context"
	"encoding/json/v2"
	"errors"
	"jobradar/tracker-service/internal/domain"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeApplicationStore struct {
    createErr error
    listErr error
    apps []domain.Application
}

func (f *fakeApplicationStore) Create(ctx context.Context, app *domain.Application) error {
    if f.createErr != nil {
        return f.createErr
    }

    app.Id = int64(len(f.apps) + 1)
    f.apps = append(f.apps, *app)

    return nil
}

func (f *fakeApplicationStore) List(ctx context.Context) ([]domain.Application, error) {
    if f.listErr != nil {
        return nil, f.listErr
    }

    return f.apps, nil
}

func TestCreateApplication_ValidBody(t *testing.T) {
	body := strings.NewReader(`{"company":"Acme","role":"Backend Engineer","url":"https://acme.com/jobs/1"}`)
	req := httptest.NewRequest(http.MethodPost, "/applications", body)
	rec := httptest.NewRecorder()

	handler := NewApplicationHandler(
		slog.New(slog.DiscardHandler),
		&fakeApplicationStore{
			createErr: nil,
			listErr: nil,
			apps: []domain.Application{},
		},
	)
	handler.CreateApplication(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var data domain.Application
	err := json.Unmarshal((rec.Body.Bytes()), &data)

	require.NoError(t, err)

	assert.Equal(t, "Acme", data.Company)
	assert.Equal(t, "Backend Engineer", data.Role)
	assert.Equal(t, "https://acme.com/jobs/1", data.OfferUrl)	
}

func TestCreateApplication_InvalidBody(t *testing.T) {
	body := strings.NewReader(`{company":"Acme","role":"Backend Engineer","url":"https://acme.com/jobs/1"}`)
	req := httptest.NewRequest(http.MethodPost, "/applications", body)
	rec := httptest.NewRecorder()

	handler := NewApplicationHandler(
		slog.New(slog.DiscardHandler),
		&fakeApplicationStore{
			createErr: nil,
			listErr: nil,
			apps: []domain.Application{},
		},
	)
	handler.CreateApplication(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var data domain.Application
	err := json.Unmarshal((rec.Body.Bytes()), &data)

	require.Error(t, err)
}

func TestCreateApplication_EmptyCompany(t *testing.T) {
	body := strings.NewReader(`{"company":"","role":"Backend Engineer","url":"https://acme.com/jobs/1"}`)
	req := httptest.NewRequest(http.MethodPost, "/applications", body)
	rec := httptest.NewRecorder()

	handler := NewApplicationHandler(
		slog.New(slog.DiscardHandler),
		&fakeApplicationStore{
			createErr: nil,
			listErr: nil,
			apps: []domain.Application{},
		},
	)
	handler.CreateApplication(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	var data domain.Application
	err := json.Unmarshal((rec.Body.Bytes()), &data)

	require.Error(t, err)
}

func TestCreateApplication_CreateError(t *testing.T) {
	body := strings.NewReader(`{"company":"Acme","role":"Backend Engineer","url":"https://acme.com/jobs/1"}`)
	req := httptest.NewRequest(http.MethodPost, "/applications", body)
	rec := httptest.NewRecorder()

	handler := NewApplicationHandler(
		slog.New(slog.DiscardHandler),
		&fakeApplicationStore{
			createErr: errors.New("something went wrong"),
			listErr: nil,
			apps: []domain.Application{},
		},
	)
	handler.CreateApplication(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListApplications_Valid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/applications", nil)
	rec := httptest.NewRecorder()

	now := time.Now()

	handler := NewApplicationHandler(
		slog.New(slog.DiscardHandler),
		&fakeApplicationStore{
			createErr: nil,
			listErr: nil,
			apps: []domain.Application{
				{
					Id: 1,
					Company: "Acme",
					Role: "Backend Engineer",
					OfferUrl: "https://acme.com/jobs/1",
					Status: domain.StatusApplied,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					Id: 2,
					Company: "Google",
					Role: "Go Backend Engineer",
					OfferUrl: "https://acme.com/jobs/2",
					Status: domain.StatusInterview,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
	)
	handler.ListApplications(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var applications = []domain.Application{}

	err := json.Unmarshal(rec.Body.Bytes(), &applications)

	require.NoError(t, err)

	assert.Len(t, applications, 2)
	assert.Equal(t, "Acme", applications[0].Company)
	assert.Equal(t, "Backend Engineer", applications[0].Role)
	assert.Equal(t, "https://acme.com/jobs/1", applications[0].OfferUrl)
}
