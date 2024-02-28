package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/arayofcode/tokeniser/models"
)

func (rc *routerConfig) setupRoutes() {
	rc.router.GET("/ping", rc.handlePing)
	rc.router.POST("/tokenise", rc.handleTokenise)
	rc.router.POST("/detokenise", rc.handleDetokenise)
	rc.router.GET("/tokenise", rc.handleTokeniseForm)
	rc.router.GET("/detokenise", rc.handleDetokeniseForm)
	rc.router.GET("/dashboard", rc.handleDashboard)
	rc.router.GET("/all", rc.handleAll)
	rc.router.POST("/unmask", rc.handleUnmask)
}

func (rc *routerConfig) StartAPI() {
	rc.router.LoadHTMLGlob("router/templates/*")
	rc.router.Static("/static", "router/static")
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

func (rc *routerConfig) handleTokeniseForm(c *gin.Context) {
	c.HTML(http.StatusOK, "tokenise.html", nil)
}

func (rc *routerConfig) handleDetokeniseForm(c *gin.Context) {
	c.HTML(http.StatusOK, "detokenise.html", nil)
}

func (rc *routerConfig) handleDashboard(c *gin.Context) {
	ctx := context.Background()

	var payload struct {
		Mask bool `json:"mask"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		payload.Mask = true
	}

	data, err := rc.handler.GetAllCards(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Some finding all cards: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if payload.Mask {
		for i := range data {
			cardNumberDecryptedByte, _ := rc.cipher.Decrypt(data[i].CardNumberEncrypted)
			expiryDateDecryptedByte, _ := rc.cipher.Decrypt(data[i].ExpiryDateEncrypted)
			data[i].CardNumber = maskLeft(string(cardNumberDecryptedByte))
			expiryDate := string(expiryDateDecryptedByte)
			data[i].ExpiryDate = expiryDate[:2] + "/" + expiryDate[2:]
			data[i].CardNumberEncrypted = nil
			data[i].ExpiryDateEncrypted = nil
		}
	} else {
		for i := range data {
			cardNumberDecryptedByte, _ := rc.cipher.Decrypt(data[i].CardNumberEncrypted)
			expiryDateDecryptedByte, _ := rc.cipher.Decrypt(data[i].ExpiryDateEncrypted)
			data[i].CardNumber = string(cardNumberDecryptedByte)
			expiryDate := string(expiryDateDecryptedByte)
			data[i].ExpiryDate = expiryDate[:2] + "/" + expiryDate[2:]
			data[i].CardNumberEncrypted = nil
			data[i].ExpiryDateEncrypted = nil
		}
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{"data": data})
}

func (rc *routerConfig) handleAll(c *gin.Context) {
	ctx := context.Background()

	var payload struct {
		Mask bool `json:"mask" binding:"required"`
	}

	if err := c.ShouldBind(&payload); err != nil {
		payload.Mask = true
	}

	data, err := rc.handler.GetAllCards(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Some finding all cards: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if payload.Mask {
		for i := range data {
			cardNumberDecryptedByte, _ := rc.cipher.Decrypt(data[i].CardNumberEncrypted)
			expiryDateDecryptedByte, _ := rc.cipher.Decrypt(data[i].ExpiryDateEncrypted)
			data[i].CardNumber = maskLeft(string(cardNumberDecryptedByte))
			expiryDate := string(expiryDateDecryptedByte)
			data[i].ExpiryDate = expiryDate[:2] + "/" + expiryDate[2:]
			data[i].CardNumberEncrypted = nil
			data[i].ExpiryDateEncrypted = nil
		}
	} else {
		for i := range data {
			cardNumberDecryptedByte, _ := rc.cipher.Decrypt(data[i].CardNumberEncrypted)
			expiryDateDecryptedByte, _ := rc.cipher.Decrypt(data[i].ExpiryDateEncrypted)
			data[i].CardNumber = string(cardNumberDecryptedByte)
			expiryDate := string(expiryDateDecryptedByte)
			data[i].ExpiryDate = expiryDate[:2] + "/" + expiryDate[2:]
			data[i].CardNumberEncrypted = nil
			data[i].ExpiryDateEncrypted = nil
		}
	}

	c.JSON(http.StatusOK, data)
}

func (rc *routerConfig) handleUnmask(c *gin.Context) {
	ctx := context.Background()

	var payload struct {
		Token uuid.UUID `json:"token" binding:"required"`
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		log.Printf("Invalid Token provided: %s", payload.Token)
		return
	}

	response, err := rc.handler.HandleDetokenise(ctx, models.DetokenisePayload{Token: payload.Token})
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
	response.Card.ExpiryDate = response.Card.ExpiryDate[:2] + "/" + response.Card.ExpiryDate[2:]
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decrypting card number: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(http.StatusFound, response)
}

func maskLeft(s string) string {
	rs := []rune(s)
	for i := 0; i < len(rs)-4; i++ {
		rs[i] = 'X'
	}
	return string(rs)
}
