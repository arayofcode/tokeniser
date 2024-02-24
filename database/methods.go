package database

import (
	"context"
	"fmt"
	"github.com/arayofcode/tokeniser/models"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// func (dbConfig *databaseConfig) InsertPII(ctx context.Context, person models.Person) {
// 	conn, err := pgx.Connect(ctx, dbConfig.DB_URL)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
// 		os.Exit(1)
// 	}
// 	defer conn.Close(ctx)

// 	conn.Query(ctx, "INSERT INTO person")

// }

func (dbConfig *databaseConfig) InsertCard(ctx context.Context, card models.Card) (results models.InsertCardResult) {
	conn, err := pgx.Connect(ctx, dbConfig.DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(ctx)

	var dbcard insertCardResult
	err = conn.QueryRow(ctx, "INSERT INTO credit_cards (cardholder_name, card_number, expiry_date, card_number_encrypted, expiry_date_encrypted) VALUES ($1, $2, $3, $4, $5) RETURNING id, token, created_at, updated_at",
		card.CardHolderName,
		card.CardNumber,
		card.ExpiryDate,
		card.CardNumberEncrypted,
		card.ExpirydateEncrypted,
	).Scan(
		&dbcard.ID,
		&dbcard.Token,
		&dbcard.CreatedAt,
		&dbcard.UpdatedAt,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert card details: %v\n", err)
		return
	}

	results = models.InsertCardResult{
		ID:        dbcard.ID,
		Token:     uuid.UUID(dbcard.Token.Bytes),
		CreatedAt: dbcard.CreatedAt.Time,
		UpdatedAt: dbcard.UpdatedAt.Time,
	}
	return
}

func (dbConfig *databaseConfig) GetCardDetails(ctx context.Context, token string) (creditCard models.Card) {
	conn, err := pgx.Connect(ctx, dbConfig.DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(ctx)

	var dbcard card
	err = conn.QueryRow(ctx, "SELECT id, token, cardholder_name, card_number, expiry_date, created_at, updated_at, card_number_encrypted, expiry_date_encrypted FROM credit_cards WHERE token=$1",
		token,
	).Scan(
		&dbcard.ID,
		&dbcard.Token,
		&dbcard.CardHolderName,
		&dbcard.CardNumber,
		&dbcard.ExpiryDate,
		&dbcard.CreatedAt,
		&dbcard.UpdatedAt,
		&dbcard.CardNumberEncrypted,
		&dbcard.ExpirydateEncrypted,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving credit card details with token %s: %v\n", token, err)
	}

	creditCard.ID = dbcard.ID
	creditCard.Token = dbcard.Token.Bytes
	creditCard.CardHolderName = dbcard.CardHolderName.String
	creditCard.CardNumber = dbcard.CardNumber.String
	creditCard.CardNumberEncrypted = dbcard.CardNumberEncrypted
	creditCard.ExpiryDate = dbcard.ExpiryDate.String
	creditCard.ExpirydateEncrypted = dbcard.ExpirydateEncrypted
	creditCard.CreatedAt = dbcard.CreatedAt.Time
	creditCard.UpdatedAt = dbcard.UpdatedAt.Time
	return
}

func (dbConfig *databaseConfig) DeleteCard(ctx context.Context, token string) (deleted bool) {
	conn, err := pgx.Connect(ctx, dbConfig.DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(ctx)

	transaction, err := conn.Begin(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start transaction: %v\n", err)
		return false
	}

	defer transaction.Rollback(ctx)

	cmdTag, err := transaction.Exec(ctx, "DELETE FROM credit_cards WHERE token=$1", token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Delete Failed: %v\n", err)
		return false
	}

	if cmdTag.RowsAffected() == 0 {
		fmt.Println("No rows were deleted.")
		return false
	} else {
		fmt.Printf("%d rows were deleted.\n", cmdTag.RowsAffected())
	}

	if err := transaction.Commit(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to commit transaction: %v\n", err)
		return false
	}

	return true
}

func (dbConfig *databaseConfig) TempShowCards(ctx context.Context) (cards []models.Card) {
	conn, err := pgx.Connect(ctx, dbConfig.DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(ctx)
	rows, err := conn.Query(ctx, "SELECT id, token, cardholder_name, card_number, expiry_date, created_at, updated_at, card_number_encrypted, expiry_date_encrypted FROM credit_cards")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while querying: %v\n", err)
		return
	}

	for rows.Next() {
		var cardRow card
		err = rows.Scan(
			&cardRow.ID,
			&cardRow.Token,
			&cardRow.CardHolderName,
			&cardRow.CardNumber,
			&cardRow.ExpiryDate,
			&cardRow.CreatedAt,
			&cardRow.UpdatedAt,
			&cardRow.CardNumberEncrypted,
			&cardRow.ExpirydateEncrypted,
		)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while parsing row: %v\n", err)
			return
		}

		var creditCard = models.Card{
			ID:                  cardRow.ID,
			Token:               uuid.UUID(cardRow.Token.Bytes[:]),
			CardHolderName:      cardRow.CardHolderName.String,
			CardNumber:          cardRow.CardNumber.String,
			ExpiryDate:          cardRow.ExpiryDate.String,
			CreatedAt:           cardRow.CreatedAt.Time,
			UpdatedAt:           cardRow.UpdatedAt.Time,
			CardNumberEncrypted: cardRow.CardNumberEncrypted,
			ExpirydateEncrypted: cardRow.ExpirydateEncrypted,
		}
		cards = append(cards, creditCard)
	}
	return
}
