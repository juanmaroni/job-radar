package domain

import (
	"errors"
	"time"
)

type Job struct {
    Id int64
    Source string
    SourceId string
    Company string
    Role string
	OfferUrl string
	Location string
    SalaryMin *int32
    SalaryMax *int32
	Description string
    Score int32
	PublishedAt *time.Time
    CreatedAt time.Time
}

var ErrInvalidSource = errors.New("no source provided")
var ErrInvalidSourceId = errors.New("no source id provided")
var ErrInvalidCompany = errors.New("no company provided")
var ErrInvalidRole = errors.New("no role provided")
var ErrInvalidOfferUrl = errors.New("no offer url provided")
var ErrInvalidScore = errors.New("invalid score provided")

func NewJob(source, sourceId, company, role, offerUrl, location, description string, score int32, salaryMin *int32, salaryMax *int32, published_at *time.Time) (*Job, error) {
	if source == "" {
		return nil, ErrInvalidSource
	}

	if sourceId == "" {
		return nil, ErrInvalidSourceId
	}

	if company == "" {
		return nil, ErrInvalidCompany
	}

	if role == "" {
		return nil, ErrInvalidRole
	}

	if offerUrl == "" {
		return nil, ErrInvalidOfferUrl
	}

	if salaryMin != nil && *salaryMin < 0 {
		salaryMin = nil
	}
	
	if salaryMax != nil && *salaryMax < 0 {
		salaryMax = nil
	}

	if score < 0 {
		return nil, ErrInvalidScore
	}

	return &Job{
		Source: source,
		SourceId: sourceId,
		Company: company,
		Role: role,
		OfferUrl: offerUrl,
		Location: location,
		SalaryMin: salaryMin,
		SalaryMax: salaryMax,
		Description: description,
		Score: score,
		PublishedAt: published_at,
		CreatedAt: time.Now().UTC(),
	}, nil
}
