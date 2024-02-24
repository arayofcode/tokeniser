package models

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID                  int       `json:"id" db:"id"`
	Token               uuid.UUID `json:"token" db:"token"`
	CardHolderName      string    `json:"cardholder_name" db:"cardholder_name"`
	CardNumber          string    `json:"card_number" db:"card_number"`
	ExpiryDate          string    `json:"expiry_date" db:"expiry_date"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
	CardNumberEncrypted []byte    `json:"card_number_encrypted" db:"card_number_encrypted"`
	ExpirydateEncrypted []byte    `json:"expiry_date_encrypted" db:"expiry_date_encrypted"`
}

type InsertCardResult struct {
	ID        int       `json:"id"`
	Token     uuid.UUID `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
