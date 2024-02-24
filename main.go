package main

import (
	// "context"
	"fmt"
	"net/http"

	// "github.com/arayofcode/tokeniser/common"
	// "github.com/arayofcode/tokeniser/database"
	// "github.com/arayofcode/tokeniser/handler"
	"github.com/arayofcode/tokeniser/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var dataStore map[string]string

func HandlePing(c *gin.Context) {
	c.String(200, "pong")
}

func GenerateToken(key string, value string) (token string) {
	fmt.Println("Creating token for key", key)
	token = uuid.NewString()
	dataStore[token] = value
	return
}

// Returns token for each value generated
// Not handling idempotency as of now
func HandleTokenize(c *gin.Context) {
	if dataStore == nil {
		dataStore = make(map[string]string)
	}

	var newRequest models.Payload
	if err := c.BindJSON(&newRequest); err != nil {
		return
	}
	var responseData models.Payload
	responseData.ID = newRequest.ID
	responseData.Data = make(map[string]string)
	for key, value := range newRequest.Data {
		responseData.Data[key] = GenerateToken(key, value)
	}

	c.JSON(http.StatusCreated, responseData)
}

func Detokenize(token string) (exists bool, value string) {
	if value, exists = dataStore[token]; !exists {
		value = "invalid token"
	}
	return
}

func HandleDetokenize(c *gin.Context) {
	var newRequest models.Payload
	if err := c.BindJSON(&newRequest); err != nil {
		return
	}

	var responseData models.DeTokenizeResponse
	responseData.ID = newRequest.ID
	responseData.Data = make(map[string]models.DeTokenizeResponseData)
	for field, token := range newRequest.Data {
		exists, originalValue := Detokenize(token)
		responseData.Data[field] = models.DeTokenizeResponseData{
			Found: exists,
			Value: originalValue,
		}
	}

	c.JSON(http.StatusCreated, responseData)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", HandlePing)
	r.POST("/tokenize", HandleTokenize)
	r.POST("/detokenize", HandleDetokenize)
	return r
}

func main() {
	// ctx := context.Background()

	// db := database.DatabaseInit(common.GetDbURL())
	// newHandler := handler.NewHandler(db)
	// databaseUrl := common.GetDbURL()
	// dbconfig := database.DatabaseInit(databaseUrl)
	// results := dbconfig.TempShowCards(ctx)
	// fmt.Printf("%+v\n", results)

	r := setupRouter()
	r.Run(":8080")
}

/*
Possible Idempotency:
- Request with same payload sent again
- Same key-value pair sent again. Should we generate new tokens in this case? Think about same names of two different people

User --> API --> Handler --> Database
*/
