package models

import (
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID   int               `json:"id"`
	Data map[string]string `json:"data"`
}

type TokenisePayload struct {
	ID   int               `json:"id"`
	Card CreditCardDetails `json:"card"`
}

type DetokenisePayload struct {
	ID    int       `json:"id"`
	Token uuid.UUID `json:"token"`
}

type DetokeniseCardResponse struct {
	ID   int               `json:"id"`
	Card CreditCardDetails `json:"card"`
}

type InsertCardResult struct {
	ID        int       `json:"id"`
	Token     uuid.UUID `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TokeniseCardResponse struct {
	InsertCardResult
}

type DeTokeniseResponseData struct {
	Found bool   `json:"found"`
	Value string `json:"value"`
}

type DeTokeniseResponse struct {
	ID   int                               `json:"id"`
	Data map[string]DeTokeniseResponseData `json:"data"`
}
