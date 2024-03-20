package handler

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/arayofcode/tokeniser/models"
)

func (h *HandlerData) HandleTokenise(ctx context.Context, newPayload models.TokenisePayload) (response models.TokeniseCardResponse, err error) {
	insertCard, err := h.db.InsertCard(ctx, newPayload.Card)
	if err != nil {
		log.Error().Err(err).Str("requestID", newPayload.RequestID).Msg("Failed to insert card")
		return models.TokeniseCardResponse{}, err
	}

	response = models.TokeniseCardResponse{
		RequestID: newPayload.RequestID,
		InsertCardResult: models.InsertCardResult{
			RowID:     insertCard.RowID,
			Token:     insertCard.Token,
			CreatedAt: insertCard.CreatedAt,
			UpdatedAt: insertCard.UpdatedAt,
		},
	}

	log.Printf("Generated token %s for card request ID: %s", response.Token, response.RequestID)
	return response, err
}

func (h *HandlerData) HandleDetokenise(ctx context.Context, payload models.DetokenisePayload) (models.DetokeniseCardResponse, error) {
	cardDetails, err := h.db.GetCardDetails(ctx, payload.Token)
	if err != nil {
		log.Error().Err(err).Str("token", payload.Token.String()).Msg("Failed to retrieve card details")
		return models.DetokeniseCardResponse{}, err
	}

	response := models.DetokeniseCardResponse{
		RequestID: payload.RequestID,
		Card: models.CreditCardDetails{
			CardHolderName:      cardDetails.CardHolderName,
			CardNumberEncrypted: cardDetails.CardNumberEncrypted,
			ExpiryDateEncrypted: cardDetails.ExpiryDateEncrypted,
		},
	}

	return response, nil
}

func (h *HandlerData) GetAllCards(ctx context.Context) ([]models.CreditCardRow, error) {
	cards, err := h.db.ShowAllCards(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve all cards")
		return nil, err
	}

	log.Info().Msgf("Retrieved %d cards", len(cards))
	return cards, nil
}
