package models

import (
	"time"

	"github.com/google/uuid"
)

type TokenisePayload struct {
	RequestID string            `json:"request_id"`
	Card      CreditCardDetails `json:"card"`
}

type DetokenisePayload struct {
	RequestID string    `json:"request_id"`
	Token     uuid.UUID `json:"token"`
}

type DetokeniseCardResponse struct {
	RequestID string            `json:"request_id"`
	Card      CreditCardDetails `json:"card"`
}

type InsertCardResult struct {
	RowID     int       `json:"-"`
	RequestID string    `json:"request_id"`
	Token     uuid.UUID `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

type TokeniseCardResponse struct {
	InsertCardResult
}
