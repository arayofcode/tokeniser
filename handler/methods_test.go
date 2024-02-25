package handler_test

import (
	"context"
	"log"
	"testing"

	"github.com/arayofcode/tokeniser/common"
	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/handler"
	"github.com/arayofcode/tokeniser/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	ctx      context.Context
	dbconfig database.Database
)

func configInit() {
	ctx = context.Background()
	databaseUrl := common.GetDbURL()
	dbconfig = database.DatabaseInit(databaseUrl)
}

func TestHandlerTokeniseNew(t *testing.T) {
	configInit()
	newHandler := handler.NewHandler(dbconfig)
	newPayload := models.TokenisePayload{
		ID: 12345,
		Card: models.CreditCardDetails{
			CardHolderName: "Test David",
			CardNumber:     "4112123456780987",
			ExpiryDate:     "11/22",
		},
	}
	results := newHandler.HandleTokeniseNew(ctx, newPayload)
	assert.NoError(t, uuid.Validate(results.Token.String()))
	log.Println("Cleaning up the row")
	assert.True(t, dbconfig.DeleteCard(ctx, results.Token))
}

func TestHandleDetokeniseNew(t *testing.T) {
	configInit()
	newHandler := handler.NewHandler(dbconfig)
	newPayload := models.TokenisePayload{
		ID: 12345,
		Card: models.CreditCardDetails{
			CardHolderName: "Test David",
			CardNumber:     "4112123456780987",
			ExpiryDate:     "11/22",
		},
	}
	tokeniseResults := newHandler.HandleTokeniseNew(ctx, newPayload)
	detokeniseResults := newHandler.HandleDetokeniseNew(ctx, models.DetokenisePayload{ID: newPayload.ID, Token: tokeniseResults.Token})
	assert.Equal(t, newPayload.Card.CardHolderName, detokeniseResults.Card.CardHolderName)
	assert.Equal(t, newPayload.Card.CardNumber, detokeniseResults.Card.CardNumber)
	assert.Equal(t, newPayload.Card.ExpiryDate, detokeniseResults.Card.ExpiryDate)
}