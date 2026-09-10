package events

type JobPublishedEvent struct {
    SourceId string `json:"source_id"`
    Source   string `json:"source"`
    Company  string `json:"company"`
    Role     string `json:"role"`
    OfferUrl string `json:"offer_url"`
}
