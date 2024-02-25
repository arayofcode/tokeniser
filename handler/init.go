package handler

import (
	"context"

	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/models"
)

type Handler interface {
	HandleTokeniseNew(ctx context.Context, payload models.TokenisePayload) (response models.TokeniseCardResponse)
	HandleDetokeniseNew(ctx context.Context, payload models.DetokenisePayload) (response models.DetokeniseCardResponse)
}

type HandlerData struct {
	db database.Database
}

func NewHandler(db database.Database) Handler {
	return &HandlerData{
		db: db,
	}
}
