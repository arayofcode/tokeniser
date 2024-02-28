package handler

import (
	"context"

	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/models"
)

type Handler interface {
	HandleTokenise(ctx context.Context, payload models.TokenisePayload) (response models.TokeniseCardResponse, err error)
	HandleDetokenise(ctx context.Context, payload models.DetokenisePayload) (response models.DetokeniseCardResponse, err error)
	GetAllCards(ctx context.Context) (response []models.CreditCardRow, err error)
}

type HandlerData struct {
	db database.Database
}

func NewHandler(db database.Database) Handler {
	return &HandlerData{
		db: db,
	}
}
