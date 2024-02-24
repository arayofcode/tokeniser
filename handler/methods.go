package handler

import (
	"context"
	"log"

	"github.com/arayofcode/tokeniser/models"
	"github.com/google/uuid"
)

func (h *HandlerData) HandleTokeniseNew(ctx context.Context, newPayload models.NewPayload) (response models.TokenizeCardResponse) {
	insertCard := h.db.InsertCard(ctx, newPayload.Card)
	response.ID = newPayload.ID
	response.Token = insertCard.Token
	log.Printf("Generated token %s for card request ID: %d", response.Token, response.ID)
	return response
}

func (h *HandlerData) HandleDetokeniseNew(ctx context.Context, token uuid.UUID) (response models.CreditCardRow) {
	response = h.db.GetCardDetails(ctx, token)
	return
}
