package models

type Payload struct {
	ID   string            `json:"id"`
	Data map[string]string `json:"data"`
}

type DeTokenizeResponseData struct {
	Found bool   `json:"found"`
	Value string `json:"value"`
}

type DeTokenizeResponse struct {
	ID   string `json:"id"`
	Data map[string]DeTokenizeResponseData `json:"data"`
}
