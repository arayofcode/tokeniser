package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/arayofcode/tokeniser/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (dbConfig *databaseConfig) InsertCard(ctx context.Context, card models.CreditCardDetails) (results models.InsertCardResult, err error) {
	log.Println("Inserting card details")
	conn, err := pgx.Connect(ctx, dbConfig.DB_URL)
	if err != nil {
		return results, fmt.Errorf("database connection failed: %v", err)
	}
	defer conn.Close(ctx)

	var dbcard insertCardResult
	err = conn.QueryRow(ctx, "INSERT INTO credit_cards (cardholder_name, card_number_encrypted, expiry_date_encrypted) VALUES ($1, $2, $3) RETURNING id, token, created_at, updated_at",
		card.CardHolderName,
		card.CardNumberEncrypted,
		card.ExpiryDateEncrypted,
	).Scan(
		&dbcard.ID,
		&dbcard.Token,
		&dbcard.CreatedAt,
		&dbcard.UpdatedAt,
	)

	if err != nil {
		return results, fmt.Errorf("card insertion in database failed: %v", err)
	}

	results = models.InsertCardResult{
		RowID:     dbcard.ID,
		Token:     uuid.UUID(dbcard.Token.Bytes),
		CreatedAt: dbcard.CreatedAt.Time,
		UpdatedAt: dbcard.UpdatedAt.Time,
	}

	log.Printf("Inserted successfully at timestamp %s", results.CreatedAt.String())
	return
}

func (dbConfig *databaseConfig) GetCardDetails(ctx context.Context, token uuid.UUID) (creditCard models.CreditCardRow, err error) {
	log.Printf("Searching for token: %s", token)
	conn, err := pgx.Connect(ctx, dbConfig.DB_URL)
	if err != nil {
		return creditCard, fmt.Errorf("database connection failed: %v", err)
	}
	defer conn.Close(ctx)

	var dbcard card
	err = conn.QueryRow(ctx, "SELECT id, token, cardholder_name, created_at, updated_at, card_number_encrypted, expiry_date_encrypted FROM credit_cards WHERE token=$1",
		token,
	).Scan(
		&dbcard.RowID,
		&dbcard.Token,
		&dbcard.CardHolderName,
		&dbcard.CreatedAt,
		&dbcard.UpdatedAt,
		&dbcard.CardNumberEncrypted,
		&dbcard.ExpiryDateEncrypted,
	)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return creditCard, fmt.Errorf("failed GetCreditCardDetails with token %s: %v", token, err)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("No card with such token")
		return
	}

	creditCard.RowID = dbcard.RowID
	creditCard.Token = dbcard.Token.Bytes
	creditCard.CardHolderName = dbcard.CardHolderName.String
	creditCard.CardNumberEncrypted = dbcard.CardNumberEncrypted
	creditCard.ExpirydateEncrypted = dbcard.ExpiryDateEncrypted
	creditCard.CreatedAt = dbcard.CreatedAt.Time
	creditCard.UpdatedAt = dbcard.UpdatedAt.Time

	log.Println("Found! Returning details")
	return
}

func (dbConfig *databaseConfig) DeleteCard(ctx context.Context, token uuid.UUID) (deleted bool) {
	log.Printf("Deleting card with token: %s", token)
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
		log.Println("No card found with given token")
		return false
	} else {
		log.Printf("%d rows were deleted.", cmdTag.RowsAffected())
	}

	if err := transaction.Commit(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to commit transaction: %v\n", err)
		return false
	}

	return true
}

func (dbConfig *databaseConfig) TempShowCards(ctx context.Context) (cards []models.CreditCardRow, err error) {
	log.Printf("Finding all cards in DB")

	conn, err := pgx.Connect(ctx, dbConfig.DB_URL)
	if err != nil {
		return cards, fmt.Errorf("database connection failed: %v", err)
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, "SELECT id, token, cardholder_name, card_number, expiry_date, created_at, updated_at, card_number_encrypted, expiry_date_encrypted FROM credit_cards")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while querying: %v\n", err)
		return cards, fmt.Errorf("ShowAllCards query failed: %v", err)
	}

	for rows.Next() {
		var cardRow card
		err = rows.Scan(
			&cardRow.RowID,
			&cardRow.Token,
			&cardRow.CardHolderName,
			&cardRow.CardNumber,
			&cardRow.ExpiryDate,
			&cardRow.CreatedAt,
			&cardRow.UpdatedAt,
			&cardRow.CardNumberEncrypted,
			&cardRow.ExpiryDateEncrypted,
		)

		if err != nil {
			return cards, fmt.Errorf("parsing results row failed: %v", err)
		}

		var creditCard = models.CreditCardRow{
			RowID:               cardRow.RowID,
			Token:               uuid.UUID(cardRow.Token.Bytes),
			CardHolderName:      cardRow.CardHolderName.String,
			CardNumber:          cardRow.CardNumber.String,
			ExpiryDate:          cardRow.ExpiryDate.String,
			CreatedAt:           cardRow.CreatedAt.Time,
			UpdatedAt:           cardRow.UpdatedAt.Time,
			CardNumberEncrypted: cardRow.CardNumberEncrypted,
			ExpirydateEncrypted: cardRow.ExpiryDateEncrypted,
		}
		cards = append(cards, creditCard)
	}
	log.Printf("Found %d cards", len(cards))
	return
}
