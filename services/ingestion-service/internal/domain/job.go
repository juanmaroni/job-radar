package domain

import (
	"errors"
	"time"
)

var ErrInvalidCompany = errors.New("no company provided") 
var ErrInvalidRole= errors.New("no role provided")
var ErrInvalidOfferUrl= errors.New("no offer URL provided")

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
    PublishedAt *time.Time
    CreatedAt time.Time
}

func NewJob(source, sourceId, company, role, offerUrl, location, description string, salaryMin, salaryMax *int32, publishedAt *time.Time) (*Job, error) {
	if company == "" {
		return nil, ErrInvalidCompany
	} else if role == "" {
		return nil, ErrInvalidRole
	} else if offerUrl == "" {
		return nil, ErrInvalidOfferUrl
	}
	
	if salaryMin != nil && *salaryMin < 0 {
		salaryMin = nil
	}
	
	if salaryMax != nil && *salaryMax < 0 {
		salaryMax = nil
	}

	now := time.Now().UTC()

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
		PublishedAt: publishedAt,
		CreatedAt: now,
	}, nil
}
