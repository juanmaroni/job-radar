package domain

import (
	"errors"
	"time"
)

type Job struct {
    Id       int64
    Source   string
    SourceId string
    Company  string
    Role     string
	OfferUrl string
    Score    int32
    CreatedAt time.Time
}

var ErrInvalidSource = errors.New("no source provided")
var ErrInvalidSourceId = errors.New("no source id provided")
var ErrInvalidCompany = errors.New("no company provided")
var ErrInvalidRole = errors.New("no role provided")
var ErrInvalidOfferUrl = errors.New("no offer url provided")
var ErrInvalidScore = errors.New("invalid score provided")

func NewJob(source, sourceId, company, role, offerUrl string, score int32) (*Job, error) {
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

	if score < 0 {
		return nil, ErrInvalidScore
	}

	return &Job{
		Source: source,
		SourceId: sourceId,
		Company: company,
		Role: role,
		OfferUrl: offerUrl,
		Score: score,
		CreatedAt: time.Now().UTC(),
	}, nil
}
