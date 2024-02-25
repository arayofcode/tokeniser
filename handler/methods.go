package handler

import (
	"context"
	"log"

	"github.com/arayofcode/tokeniser/models"
)

func (h *HandlerData) HandleTokeniseNew(ctx context.Context, newPayload models.TokenisePayload) (response models.TokeniseCardResponse) {
	insertCard := h.db.InsertCard(ctx, newPayload.Card)
	response.ID = newPayload.ID
	response.Token = insertCard.Token
	log.Printf("Generated token %s for card request ID: %d", response.Token, response.ID)
	return response
}

func (h *HandlerData) HandleDetokeniseNew(ctx context.Context, payload models.DetokenisePayload) (response models.DetokeniseCardResponse) {
	cardDetails := h.db.GetCardDetails(ctx, payload.Token)
	response.ID = payload.ID
	response.Card.CardHolderName = cardDetails.CardHolderName
	response.Card.CardNumber = cardDetails.CardNumber
	response.Card.ExpiryDate = cardDetails.ExpiryDate
	return
}
