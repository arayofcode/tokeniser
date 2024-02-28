package handler

import (
	"context"
	"log"

	"github.com/arayofcode/tokeniser/models"
)

func (h *HandlerData) HandleTokenise(ctx context.Context, newPayload models.TokenisePayload) (response models.TokeniseCardResponse, err error) {
	insertCard, err := h.db.InsertCard(ctx, newPayload.Card)
	response.RequestID = newPayload.RequestID
	response.RowID = insertCard.RowID
	response.Token = insertCard.Token
	response.CreatedAt = insertCard.CreatedAt
	response.UpdatedAt = insertCard.UpdatedAt
	log.Printf("Generated token %s for card request ID: %s", response.Token, response.RequestID)
	return response, err
}

func (h *HandlerData) HandleDetokenise(ctx context.Context, payload models.DetokenisePayload) (response models.DetokeniseCardResponse, err error) {
	cardDetails, err := h.db.GetCardDetails(ctx, payload.Token)
	response.RequestID = payload.RequestID
	response.Card.CardNumberEncrypted = cardDetails.CardNumberEncrypted
	response.Card.ExpiryDateEncrypted = cardDetails.ExpirydateEncrypted
	response.Card.CardHolderName = cardDetails.CardHolderName
	return response, err
}
