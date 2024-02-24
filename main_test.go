package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	// "github.com/arayofcode/tokeniser/models"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestGenerateToken(t *testing.T) {
	pairs := []string{
		"341",
		"12342",
		"3o24u",
		"",
	}

	// Initialise relevant global variables
	dataStore = make(map[string]string)
	
	for count, val := range pairs {
		assert.Equal(t, fmt.Sprint(count), GenerateToken("", val))
	}
	
	fakeDataStore := make(map[string]string)
	fakeDataStore[fmt.Sprint(0)] = "341"
	fakeDataStore[fmt.Sprint(1)] = "12342"
	fakeDataStore[fmt.Sprint(2)] = "3o24u"
	fakeDataStore[fmt.Sprint(3)] = ""
	assert.Equal(t, fakeDataStore, dataStore)
}

func TestHandleTokenization(t *testing.T) {
	// payload := `{
	// "id": req-12345”,
	// "data": {
	// 		"field1": "t6yh4f6",
	// 		"field2": "gh67ned",
	// 		"fieldn": "bnj7ytb"
	// 	}
	// }`

}
