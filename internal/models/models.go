package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type URLItem struct {
	CorrelationID string `json:"_,omitempty"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
	IsDeleted     bool   `json:"is_deleted"`
}

type BatchRequest []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchResponse []BatchResponseItem

type DeleteURLsRequest []string
