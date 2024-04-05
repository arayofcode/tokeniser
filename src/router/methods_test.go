package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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
	testRouter *gin.Engine
	db         database.Database
	h          handler.Handler
	c          cipher.Cipher
)

func configInit() {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := validate.RegisterValidation("expirydate", common.ExpiryDateMMYY)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to register 'expirydate' validation")
		}

		err = validate.RegisterValidation("notallzero", common.NotAllZero)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to register 'notallzero' validation")
		}
	}
	databaseUrl := common.GetDbURL()
	db = database.DatabaseInit(databaseUrl)
	h = handler.NewHandler(db)
	c = cipher.Init(common.GetPassphrase())
	testRouter = setupRouter()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	h := handler.NewHandler(db)
	routerConfig := &routerConfig{router: r, handler: h, cipher: c}
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
			"card_number": "4222222222222",
			"expiry_date": "1225"
		}
	}
	`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/tokens", bytes.NewBufferString(payload))
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
			"card_number": "378282246310005",
			"expiry_date": "1224"
		}
	}
	`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/tokens", bytes.NewBufferString(payload))
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
	req, _ = http.NewRequest("POST", "/v1/tokens/detokenise", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	var result models.DetokeniseCardResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	log.Printf("%+v", result)

	assert.NoError(t, err)
	assert.Equal(t, "Test_name", result.Card.CardHolderName)
	assert.Equal(t, "378282246310005", result.Card.CardNumber)
	assert.Equal(t, "12/24", result.Card.ExpiryDate)

	db.DeleteCard(context.TODO(), createResult.Token)
}
