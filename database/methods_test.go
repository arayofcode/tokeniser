package database_test

import (
	"context"
	"testing"

	"github.com/arayofcode/tokeniser/common"
	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/models"

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

func TestInsert(t *testing.T) {
	configInit()
	results := dbconfig.InsertCard(ctx, models.Card{
		CardHolderName:      "Test",
		CardNumber:          "4111098712348484",
		ExpiryDate:          "11/22",
		ExpirydateEncrypted: []byte{},
		CardNumberEncrypted: []byte{},
	})
	assert.NotNil(t, results)
	common.PrettyPrint(results)
}

func TestRetrieve(t *testing.T) {
	configInit()
	card := models.Card{
		CardHolderName:      "Test",
		CardNumber:          "4111098712348484",
		ExpiryDate:          "11/22",
	}
	results := dbconfig.InsertCard(ctx, card)
	searchResults := dbconfig.GetCardDetails(ctx, results.Token.String())
	common.PrettyPrint(searchResults)
	assert.NotNil(t, results.CreatedAt)
	assert.NotNil(t, results.UpdatedAt)
	assert.NotNil(t, results.ID)
	assert.Equal(t, card.CardHolderName, searchResults.CardHolderName)
	assert.Equal(t, card.CardNumber, searchResults.CardNumber)
	assert.Equal(t, card.ExpiryDate, searchResults.ExpiryDate)
	assert.True(t, dbconfig.DeleteCard(ctx, results.Token.String()))
	searchResults = dbconfig.GetCardDetails(ctx, results.Token.String())
	assert.Zero(t, searchResults.ID)
}