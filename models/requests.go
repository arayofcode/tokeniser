package models

import (
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID   int               `json:"id"`
	Data map[string]string `json:"data"`
}

type NewPayload struct {
	ID   int               `json:"id"`
	Card CreditCardDetails `json:"card"`
}

type InsertCardResult struct {
	ID        int       `json:"id"`
	Token     uuid.UUID `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TokenizeCardResponse struct {
	InsertCardResult
}

type DeTokenizeResponseData struct {
	Found bool   `json:"found"`
	Value string `json:"value"`
}

type DeTokenizeResponse struct {
	ID   int                               `json:"id"`
	Data map[string]DeTokenizeResponseData `json:"data"`
}
