package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/arayofcode/tokeniser/src/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (dbConfig *databaseConfig) InsertCard(ctx context.Context, card models.CreditCardDetails) (results models.InsertCardResult, err error) {
	log.Info().Msg("Inserting card details")
	conn, err := pgx.Connect(ctx, dbConfig.DB_URL)
	if err != nil {
		log.Error().Err(err).Msg("Database connection failed")
		return results, err
	}
	defer conn.Close(ctx)

	var dbcard insertCardResult
	err = conn.QueryRow(ctx, `
		INSERT INTO credit_cards (cardholder_name, card_number_encrypted, expiry_date_encrypted) 
		VALUES ($1, $2, $3) 
		RETURNING id, token, created_at, updated_at`,
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
		log.Error().Err(err).Msg("Failed to insert card into database")
		return results, err
	}

	results = models.InsertCardResult{
		RowID:     dbcard.ID,
		Token:     uuid.UUID(dbcard.Token.Bytes),
		CreatedAt: dbcard.CreatedAt.Time,
		UpdatedAt: dbcard.UpdatedAt.Time,
	}

	log.Info().
		Str("token", results.Token.String()).
		Time("createdAt", results.CreatedAt).
		Msg("Card inserted successfully")
	return
}

func (dbConfig *databaseConfig) GetCard(ctx context.Context, token uuid.UUID) (creditCard models.CreditCardRow, err error) {
	log.Info().Str("token", token.String()).Msg("Retrieving card details")
	conn, err := pgx.Connect(ctx, dbConfig.DB_URL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return creditCard, err
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

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info().Str("token", token.String()).Msg("No card found with the given token")
			err = nil
		} else {
			log.Error().Err(err).Str("token", token.String()).Msg("Failed to retrieve card details")
		}
		return creditCard, err
	}

	creditCard = models.CreditCardRow{
		RowID:               dbcard.RowID,
		Token:               dbcard.Token.Bytes,
		CardHolderName:      dbcard.CardHolderName.String,
		CardNumberEncrypted: dbcard.CardNumberEncrypted,
		ExpiryDateEncrypted: dbcard.ExpiryDateEncrypted,
		CreatedAt:           dbcard.CreatedAt.Time,
		UpdatedAt:           dbcard.UpdatedAt.Time,
	}

	log.Info().Str("token", token.String()).Msg("Card details found!")
	return
}

func (dbConfig *databaseConfig) DeleteCard(ctx context.Context, token uuid.UUID) (deleted bool) {
	log.Info().Str("token", token.String()).Msg("Attempting to delete card")
	conn, err := pgx.Connect(ctx, dbConfig.DB_URL)
	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database for deletion")
		return false
	}
	defer conn.Close(ctx)

	transaction, err := conn.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start transaction")
		return false
	}
	defer func() {
		if err := transaction.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Error().Err(err).Msg("Failed to rollback transaction")
		}
	}()

	cmdTag, err := transaction.Exec(ctx, "DELETE FROM credit_cards WHERE token=$1", token)
	if err != nil {
		log.Error().Err(err).Str("token", token.String()).Msg("Failed to delete card from database within transaction")
		return false
	}

	if cmdTag.RowsAffected() == 0 {
		log.Info().Str("token", token.String()).Msg("No card found with the given token for deletion")
		return false
	}

	if err := transaction.Commit(ctx); err != nil {
		log.Error().Err(err).Str("token", token.String()).Msg("Failed to commit transaction for card deletion")
		return false
	}

	log.Info().Str("token", token.String()).Msgf("%d row(s) were deleted.", cmdTag.RowsAffected())
	return true
}

func (dbConfig *databaseConfig) ShowAllCards(ctx context.Context) (cards []models.CreditCardRow, err error) {
	log.Printf("Finding all cards in DB")

	conn, err := pgx.Connect(ctx, dbConfig.DB_URL)
	if err != nil {
		return cards, fmt.Errorf("database connection failed: %v", err)
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, "SELECT token, cardholder_name, created_at, updated_at, card_number_encrypted, expiry_date_encrypted FROM credit_cards")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while querying: %v\n", err)
		return cards, fmt.Errorf("ShowAllCards query failed: %v", err)
	}

	for rows.Next() {
		var cardRow card
		err = rows.Scan(
			&cardRow.Token,
			&cardRow.CardHolderName,
			&cardRow.CreatedAt,
			&cardRow.UpdatedAt,
			&cardRow.CardNumberEncrypted,
			&cardRow.ExpiryDateEncrypted,
		)

		if err != nil {
			return cards, fmt.Errorf("parsing results row failed: %v", err)
		}

		var creditCard = models.CreditCardRow{
			Token:               uuid.UUID(cardRow.Token.Bytes),
			CardHolderName:      cardRow.CardHolderName.String,
			CreatedAt:           cardRow.CreatedAt.Time,
			UpdatedAt:           cardRow.UpdatedAt.Time,
			CardNumberEncrypted: cardRow.CardNumberEncrypted,
			ExpiryDateEncrypted: cardRow.ExpiryDateEncrypted,
		}
		cards = append(cards, creditCard)
	}
	log.Printf("Found %d cards", len(cards))
	return
}
