package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/arayofcode/tokeniser/src/models"
)

type databaseConfig struct {
	dbPool *pgxpool.Pool
}

type Database interface {
	InsertCard(ctx context.Context, card models.CreditCardDetails) (insertCardResult models.InsertCardResult, err error)
	GetCard(ctx context.Context, token uuid.UUID) (card models.CreditCardRow, err error)
	DeleteCard(ctx context.Context, token uuid.UUID) (deleted bool)
	ShowAllCards(ctx context.Context) (cards []models.CreditCardRow, err error)
}

func DatabaseInit(DB_URL string) Database {
	if DB_URL == "" {
		log.Fatal().Msg("DB URL missing")
	}

	pool, err := pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	log.Debug().Int32("max_idle_connections", pool.Stat().MaxConns()).Msg("Connection established")

	return &databaseConfig{
		dbPool: pool,
	}
}
