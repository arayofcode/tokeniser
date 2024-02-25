package models

import (
	"time"

	"github.com/google/uuid"
)

type CardInternalData struct {
	RowID               int       `json:"id,omitempty" db:"id"`
	Token               uuid.UUID `json:"token,omitempty" db:"token"`
	CreatedAt           time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at,omitempty" db:"updated_at"`
	CardNumberEncrypted []byte    `json:"card_number_encrypted" db:"card_number_encrypted"`
	ExpirydateEncrypted []byte    `json:"expiry_date_encrypted" db:"expiry_date_encrypted"`
}

type CreditCardDetails struct {
	CardHolderName string `json:"cardholder_name" binding:"required" validate:"nonzero"`
	CardNumber     string `json:"card_number" binding:"required,credit_card" validate:"notallzero"`
	ExpiryDate     string `json:"expiry_date" binding:"required" validate:"expiry_date,nonzero"`
}

type CreditCardRow struct {
	RowID               int       `json:"id,omitempty" db:"id"`
	Token               uuid.UUID `json:"token,omitempty" db:"token"`
	CreatedAt           time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at,omitempty" db:"updated_at"`
	CardHolderName      string    `json:"cardholder_name"`
	CardNumber          string    `json:"card_number"`
	ExpiryDate          string    `json:"expiry_date"`
	CardNumberEncrypted []byte    `json:"card_number_encrypted" db:"card_number_encrypted"`
	ExpirydateEncrypted []byte    `json:"expiry_date_encrypted" db:"expiry_date_encrypted"`
}

// For testing
type CreditCardRowNew struct {
	CardInternalData
	CreditCardDetails
}
