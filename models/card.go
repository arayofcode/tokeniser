package models

import (
	"time"

	"github.com/google/uuid"
)

type CreditCardDetails struct {
	CardHolderName      string `json:"cardholder_name"`
	CardNumber          string `json:"card_number"`
	ExpiryDate          string `json:"expiry_date"`
	CardNumberEncrypted []byte `json:"card_number_encrypted" db:"card_number_encrypted"`
	ExpirydateEncrypted []byte `json:"expiry_date_encrypted" db:"expiry_date_encrypted"`
}

type CreditCardRow struct {
	ID        int       `json:"id,omitempty" db:"id"`
	Token     uuid.UUID `json:"token,omitempty" db:"token"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
	CardHolderName      string `json:"cardholder_name"`
	CardNumber          string `json:"card_number"`
	ExpiryDate          string `json:"expiry_date"`
	CardNumberEncrypted []byte `json:"card_number_encrypted" db:"card_number_encrypted"`
	ExpirydateEncrypted []byte `json:"expiry_date_encrypted" db:"expiry_date_encrypted"`
}
