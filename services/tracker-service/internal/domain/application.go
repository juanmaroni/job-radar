package domain

import (
	"errors"
	"time"
)

var ErrInvalidCompany = errors.New("no company provided") 
var ErrInvalidRole= errors.New("no role provided") 

type Status string

func (s Status) Valid() bool {
	return s == StatusApplied || s == StatusInterview || s == StatusOffer || s == StatusRejected
}

const (
	StatusApplied Status = "applied"
	StatusInterview Status = "interview"
	StatusOffer Status = "offer"
	StatusRejected Status = "rejected"
)

type Application struct {
	Id int64
	Company string
	Role string
	OfferUrl string
	Status Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewApplication(company, role, url string) (*Application, error) {
	if company == "" {
		return nil, ErrInvalidCompany
	} else if role == "" {
		return nil, ErrInvalidRole
	}

	now := time.Now().UTC()

	return &Application {
		Company: company,
		Role: role,
		OfferUrl: url,
		Status: StatusApplied,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
