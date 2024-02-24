package handler

import (
	"context"

	"github.com/google/uuid"
	
	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/models"
)

type Handler interface {
	HandleTokeniseNew(ctx context.Context, newPayload models.NewPayload) (response models.TokenizeCardResponse)
	HandleDetokeniseNew(ctx context.Context, token uuid.UUID) (response models.CreditCardRow)
}

type HandlerData struct {
	db  database.Database
}

func NewHandler(db database.Database) Handler {
	return &HandlerData{
		db:  db,
	}
}
