package database_test

import (
	"context"
	"testing"

	"github.com/arayofcode/tokeniser/cipher"
	"github.com/arayofcode/tokeniser/common"
	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/models"

	"github.com/stretchr/testify/assert"
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

func TestInsert(t *testing.T) {
	configInit()

	card := models.CreditCardDetails{
		CardHolderName: "Test",
		CardNumber:     "4222222222222",
		ExpiryDate:     "11/26",
	}
	card.CardNumberEncrypted, _ = secure.Encrypt([]byte(card.CardNumber))
	card.ExpiryDateEncrypted, _ = secure.Encrypt([]byte(card.ExpiryDateEncrypted))

	results, _ := dbconfig.InsertCard(ctx, card)

	assert.NotNil(t, results)
	dbconfig.DeleteCard(ctx, results.Token)
}

func TestRetrieve(t *testing.T) {
	configInit()
	card := models.CreditCardDetails{
		CardHolderName: "Test",
		CardNumber:     "4222222222222",
		ExpiryDate:     "11/24",
	}
	card.CardNumberEncrypted, _ = secure.Encrypt([]byte(card.CardNumber))
	card.ExpiryDateEncrypted, _ = secure.Encrypt([]byte(card.ExpiryDateEncrypted))
	results, _ := dbconfig.InsertCard(ctx, card)
	searchResults, _ := dbconfig.GetCardDetails(ctx, results.Token)
	assert.NotNil(t, results.CreatedAt)
	assert.NotNil(t, results.UpdatedAt)
	assert.NotNil(t, results.RowID)
	assert.Equal(t, card.CardHolderName, searchResults.CardHolderName)
	assert.Equal(t, card.CardNumberEncrypted, searchResults.CardNumberEncrypted)
	assert.Equal(t, card.ExpiryDateEncrypted, searchResults.ExpiryDateEncrypted)

	// No need to test delete because it's already happening
	assert.True(t, dbconfig.DeleteCard(ctx, results.Token))
	searchResults, _ = dbconfig.GetCardDetails(ctx, results.Token)
	assert.Zero(t, searchResults.RowID)
}
