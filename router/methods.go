package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/arayofcode/tokeniser/models"
)

func (rc *routerConfig) setupRoutes() {
	rc.router.GET("/ping", rc.handlePing)
	rc.router.POST("/tokenise", rc.handleTokenise)
	rc.router.POST("/detokenise", rc.handleDetokenise)
}

func (rc *routerConfig) StartAPI() {
	rc.setupRoutes()
	rc.router.Run(":8080")
}

func (rc *routerConfig) handlePing(c *gin.Context) {
	c.String(200, "pong")
}

func (rc *routerConfig) handleTokenise(c *gin.Context) {
	ctx := context.Background()

	var payload models.TokenisePayload
	if err := c.BindJSON(&payload); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing payload: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload."})
		return
	}

	log.Println("Received /tokenise from request id: " + payload.RequestID)

	log.Println("Encrypting card details:")

	var err error
	payload.Card.CardNumberEncrypted, err = rc.cipher.Encrypt([]byte(payload.Card.CardNumber))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encrypting credit card number: %s\n", err)
	}

	payload.Card.ExpiryDateEncrypted, err = rc.cipher.Encrypt([]byte(payload.Card.ExpiryDate))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encrypting expiry date: %s\n", err)
	}

	log.Print("Successfully encrypted. Tokenising it")
	response, err := rc.handler.HandleTokenise(ctx, payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error tokenising: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (rc *routerConfig) handleDetokenise(c *gin.Context) {
	ctx := context.Background()

	var payload models.DetokenisePayload
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload."})
		fmt.Fprintf(os.Stderr, "Error parsing payload: %s\n", err)
		return
	}

	log.Println("Received /detokenise from request id: " + payload.RequestID)

	response, err := rc.handler.HandleDetokenise(ctx, payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	log.Println("Received encrypted card details. Decrypting")
	cardNumberByte, err := rc.cipher.Decrypt(response.Card.CardNumberEncrypted)
	response.Card.CardNumber = string(cardNumberByte)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decrypting card number: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	expiryDateByte, err := rc.cipher.Decrypt(response.Card.ExpiryDateEncrypted)
	response.Card.ExpiryDate = string(expiryDateByte)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decrypting card number: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(http.StatusFound, response)
}
