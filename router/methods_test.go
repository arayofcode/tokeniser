package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arayofcode/tokeniser/common"
	"github.com/arayofcode/tokeniser/database"
	"github.com/arayofcode/tokeniser/handler"
	"github.com/arayofcode/tokeniser/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	testRouter *gin.Engine
	db         database.Database
	h          handler.Handler
)

func configInit() {
	databaseUrl := common.GetDbURL()
	db = database.DatabaseInit(databaseUrl)
	h = handler.NewHandler(db)
	testRouter = setupRouter()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	h := handler.NewHandler(db)
	routerConfig := &routerConfig{router: r, handler: h}
	routerConfig.setupRoutes()
	return r
}

func TestHandlePing(t *testing.T) {
	configInit()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestHandleTokenise(t *testing.T) {
	configInit()

	payload := `
	{
		"request_id": "req-12345",
		"card": {
        	"cardholder_name" : "Test",
			"card_number": "4000056655665556",
			"expiry_date": "12/24"
		}
	}
	`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tokenise", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	var result models.TokeniseCardResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	err = uuid.Validate(result.Token.String())
	assert.NoError(t, err)

	db.DeleteCard(context.TODO(), result.Token)
}

func TestHandleDetokenise(t *testing.T) {
	configInit()

	// Creating mock row and generating token
	payload := `
	{
		"request_id": "req-12345",
		"card": {
        	"cardholder_name" : "Test_name",
			"card_number": "4000056655665556",
			"expiry_date": "12/24"
		}
	}
	`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tokenise", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	var createResult models.TokeniseCardResponse
	_ = json.Unmarshal(w.Body.Bytes(), &createResult)

	// Detokenising
	payload = fmt.Sprintf(`
	{
		"request_id": "req-12345",
		"token": "%s"
	}`, createResult.Token)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/detokenise", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	var result models.DetokeniseCardResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "Aryan", result.Card.CardHolderName)
	assert.Equal(t, "4000056655665556", result.Card.CardNumber)
	assert.Equal(t, "12/24", result.Card.ExpiryDate)

	db.DeleteCard(context.TODO(), createResult.Token)
}
