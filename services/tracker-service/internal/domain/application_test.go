package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApplication_ValidInput(t *testing.T) {
    app, err := NewApplication("Acme", "Backend Engineer", "https://acme.com/jobs/1", nil)

    require.NoError(t, err)
    require.NotNil(t, app)

    assert.Equal(t, "Acme", app.Company)
	assert.Equal(t, "Backend Engineer", app.Role)
	assert.Equal(t, "https://acme.com/jobs/1", app.OfferUrl)
    assert.Equal(t, StatusApplied, app.Status)
	assert.False(t, time.Time.IsZero(app.CreatedAt))
	assert.False(t, time.Time.IsZero(app.UpdatedAt))
}

func TestNewApplication_Validation(t *testing.T) {
    cases := []struct {
        name string
        company string
        role string
		matchScore *int32
        wantErr error
    } {
		{"valid application", "Acme", "Backend Engineer", nil, nil},
        {"empty company", "", "Backend Engineer", nil, ErrInvalidCompany},
        {"empty role", "Acme", "", nil, ErrInvalidRole},
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            _, err := NewApplication(c.company, c.role, "https://example.com", nil)
			assert.ErrorIs(t, err, c.wantErr)
        })
    }
}

func TestIsValid_Validation(t *testing.T) {
	var StatusFake Status = "fake"

    cases := []struct {
        name    string
		status Status
        want bool
    } {
		{"applied status", StatusApplied, true},
		{"interview status", StatusInterview, true},
		{"offer status", StatusOffer, true},
		{"rejected status", StatusRejected, true},
        {"fake status", StatusFake, false},
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, c.status.Valid())
        })
    }
}
