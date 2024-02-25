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
	results, _ := dbconfig.InsertCard(ctx, models.CreditCardDetails{
		CardHolderName: "Test",
		CardNumber:     "4111098712348484",
		ExpiryDate:     "11/22",
	})
	assert.NotNil(t, results)
	common.PrettyPrint(results)
	dbconfig.DeleteCard(ctx, results.Token)
}

func TestRetrieve(t *testing.T) {
	configInit()
	card := models.CreditCardDetails{
		CardHolderName: "Test",
		CardNumber:     "4111098712348484",
		ExpiryDate:     "11/22",
	}
	results, _ := dbconfig.InsertCard(ctx, card)
	searchResults, _ := dbconfig.GetCardDetails(ctx, results.Token)
	common.PrettyPrint(searchResults)
	assert.NotNil(t, results.CreatedAt)
	assert.NotNil(t, results.UpdatedAt)
	assert.NotNil(t, results.RowID)
	assert.Equal(t, card.CardHolderName, searchResults.CardHolderName)
	assert.Equal(t, card.CardNumber, searchResults.CardNumber)
	assert.Equal(t, card.ExpiryDate, searchResults.ExpiryDate)

	// No need to test delete because it's already happening
	assert.True(t, dbconfig.DeleteCard(ctx, results.Token))
	searchResults, _ = dbconfig.GetCardDetails(ctx, results.Token)
	assert.Zero(t, searchResults.RowID)
}
