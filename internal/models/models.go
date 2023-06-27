package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type BatchDBItem struct {
	CorrelationID string
	OriginalURL   string
	ShortURL      string
}

type BatchDB []BatchDBItem

type BatchRequest []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchResponse []BatchResponseItem
