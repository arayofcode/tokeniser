/*
Testing not yet implemented. Do this later
*/

package router

import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/arayofcode/tokeniser/models"
	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func configInit() {
	gin.SetMode(gin.TestMode)
	router = gin.Default()
	rc := &routerConfig{router: router}
	rc.setupRoutes()
}

// func TestHandlePing(t *testing.T) {
// 	configInit()
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/ping", nil)
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Equal(t, "pong", w.Body.String())
// }

// func TestHandleTokenise(t *testing.T) {
// 	configInit()
// 	payload := models.TokenisePayload{
// 		RequestID: "1234",
// 		Card: models.CreditCardDetails{
// 			CardHolderName: "Test Ray",
// 			CardNumber:     "1111-2222-3333-4444",
// 			ExpiryDate:     "11/22",
// 		},
// 	}
// 	jsonPayload, _ := json.Marshal(payload)
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/tokenise", bytes.NewBufferString(string(jsonPayload)))
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusCreated, w.Code)
// }

// func TestHandleDetokenise(t *testing.T) {
// 	configInit()
// 	payload := `{"request_id":"12345", "token":""}`
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/detokenise", bytes.NewBufferString(payload))
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusFound, w.Code)
// }
