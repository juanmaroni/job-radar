package events

type JobMatchedEvent struct {
    SourceId string `json:"source_id"`
    Source   string `json:"source"`
    Company  string `json:"company"`
    Role     string `json:"role"`
    Score    int32  `json:"score"`
}
