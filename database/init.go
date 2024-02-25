package database

import (
	"context"

	"github.com/google/uuid"
	
	"github.com/arayofcode/tokeniser/models"
)

type databaseConfig struct {
	DB_URL string
}

type Database interface {
	// InsertPII(context.Context, models.Person)
	InsertCard(ctx context.Context, card models.CreditCardDetails) (insertCardResult models.InsertCardResult, err error)
	GetCardDetails(ctx context.Context, token uuid.UUID) (card models.CreditCardRow, err error)
	DeleteCard(ctx context.Context, token uuid.UUID) (deleted bool)
	TempShowCards(ctx context.Context) (cards []models.CreditCardRow, err error)
}

func DatabaseInit(DB_URL string) Database {
	return &databaseConfig{
		DB_URL: DB_URL,
	}
}
