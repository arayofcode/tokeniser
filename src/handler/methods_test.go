package handler_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/arayofcode/tokeniser/cipher"
	"github.com/arayofcode/tokeniser/common"
	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/handler"
	"github.com/arayofcode/tokeniser/models"
)

var (
	ctx      context.Context
	dbconfig database.Database
	secure   cipher.Cipher
)

func configInit() {
	ctx = context.Background()
	databaseUrl := common.GetDbURL()
	dbconfig = database.DatabaseInit(databaseUrl)
	secure = cipher.Init(common.GetPassphrase())

}

func TestHandlerTokeniseNew(t *testing.T) {
	configInit()
	newHandler := handler.NewHandler(dbconfig)
	newPayload := models.TokenisePayload{
		RequestID: "12345",
		Card: models.CreditCardDetails{
			CardHolderName: "Test David",
			CardNumber:     "4112123456780987",
			ExpiryDate:     "11/22",
		},
	}
	newPayload.Card.CardNumberEncrypted, _ = secure.Encrypt([]byte(newPayload.Card.CardNumber))
	newPayload.Card.ExpiryDateEncrypted, _ = secure.Encrypt([]byte(newPayload.Card.ExpiryDateEncrypted))
	results, _ := newHandler.HandleTokenise(ctx, newPayload)
	assert.NoError(t, uuid.Validate(results.Token.String()))
	log.Print("Cleaning up the row")
	assert.True(t, dbconfig.DeleteCard(ctx, results.Token))
}

func TestHandleDetokeniseNew(t *testing.T) {
	configInit()
	newHandler := handler.NewHandler(dbconfig)
	newPayload := models.TokenisePayload{
		RequestID: "12345",
		Card: models.CreditCardDetails{
			CardHolderName: "Test David",
			CardNumber:     "4112123456780987",
			ExpiryDate:     "11/22",
		},
	}
	newPayload.Card.CardNumberEncrypted, _ = secure.Encrypt([]byte(newPayload.Card.CardNumber))
	newPayload.Card.ExpiryDateEncrypted, _ = secure.Encrypt([]byte(newPayload.Card.ExpiryDateEncrypted))
	tokeniseResults, _ := newHandler.HandleTokenise(ctx, newPayload)
	detokeniseResults, _ := newHandler.HandleDetokenise(ctx, models.DetokenisePayload{RequestID: newPayload.RequestID, Token: tokeniseResults.Token})
	assert.Equal(t, newPayload.Card.CardHolderName, detokeniseResults.Card.CardHolderName)
	assert.Equal(t, newPayload.Card.CardNumberEncrypted, detokeniseResults.Card.CardNumberEncrypted)
	assert.Equal(t, newPayload.Card.ExpiryDateEncrypted, detokeniseResults.Card.ExpiryDateEncrypted)
	dbconfig.DeleteCard(ctx, tokeniseResults.Token)
}
