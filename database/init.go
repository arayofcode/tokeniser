package database

import (
	"context"
	"github.com/arayofcode/tokeniser/models"
)

type databaseConfig struct {
	DB_URL string
}

type Database interface {
	// InsertPII(context.Context, models.Person)
	InsertCard(ctx context.Context, card models.Card) (insertCardResult models.InsertCardResult)
	GetCardDetails(ctx context.Context, token string) (card models.Card)
	DeleteCard(ctx context.Context, token string) (deleted bool)
	TempShowCards(ctx context.Context) []models.Card
}

func DatabaseInit(DB_URL string) Database {
	return &databaseConfig{
		DB_URL: DB_URL,
	}
}
