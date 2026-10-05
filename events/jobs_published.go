package events

import "time"

type JobPublishedEvent struct {
    SourceId string `json:"source_id"`
    Source string `json:"source"`
    Company string `json:"company"`
    Role string `json:"role"`
    OfferUrl string `json:"offer_url"`
	Location string `json:"location"`
	SalaryMin *int32 `json:"salary_min"`
	SalaryMax *int32 `json:"salary_max"`
	Description string `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
}
