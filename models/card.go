package models

import (
	"time"

	"github.com/google/uuid"
)

type CardInternalData struct {
	RowID     int       `json:"id,omitempty" db:"id"`
	Token     uuid.UUID `json:"token,omitempty" db:"token"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type CreditCardDetails struct {
	CardHolderName      string `json:"cardholder_name" binding:"required"`
	CardNumber          string `json:"card_number" binding:"required,credit_card,notallzero"`
	ExpiryDate          string `json:"expiry_date" binding:"required,expirydate"`
	CardNumberEncrypted []byte `json:"-" db:"card_number_encrypted"`
	ExpiryDateEncrypted []byte `json:"-" db:"expiry_date_encrypted"`
	
}

type CreditCardRow struct {
	RowID               int       `json:"id,omitempty" db:"id"`
	Token               uuid.UUID `json:"token,omitempty" db:"token"`
	CreatedAt           time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at,omitempty" db:"updated_at"`
	CardHolderName      string    `json:"cardholder_name"`
	CardNumber          string    `json:"card_number"`
	ExpiryDate          string    `json:"expiry_date"`
	CardNumberEncrypted []byte    `json:"card_number_encrypted,omitempty" db:"card_number_encrypted"`
	ExpiryDateEncrypted []byte    `json:"expiry_date_encrypted,omitempty" db:"expiry_date_encrypted"`
}

func (c CreditCardRow) TokenString() string {
    return c.Token.String()
}

// For testing
type CreditCardRowNew struct {
	CardInternalData
	CreditCardDetails
}
