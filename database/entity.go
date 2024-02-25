package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type card struct {
	RowID               int                `db:"id"`
	Token               pgtype.UUID        `db:"token"`
	CardHolderName      pgtype.Text        `db:"cardholder_name"`
	CardNumber          pgtype.Text        `db:"card_number"`
	ExpiryDate          pgtype.Text        `db:"expiry_date"`
	CreatedAt           pgtype.Timestamptz `db:"created_at"`
	UpdatedAt           pgtype.Timestamptz `db:"updated_at"`
	CardNumberEncrypted []byte             `db:"card_number_encrypted"`
	ExpirydateEncrypted []byte             `db:"expiry_date_encrypted"`
}

type insertCardResult struct {
	ID        int                `db:"id"`
	Token     pgtype.UUID        `db:"token"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}
