package router

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/arayofcode/tokeniser/common"
	"github.com/arayofcode/tokeniser/models"
)

func (rc *routerConfig) setupRoutes() {
	rc.router.GET("/ping", rc.handlePing)

	v1 := rc.router.Group("/v1")
	tokens := v1.Group("/tokens")
	{
		tokens.GET("", rc.handleAll)
		tokens.POST("", rc.handleTokenise)
		tokens.POST("/detokenise", rc.handleDetokenise)
		tokens.POST("/unmask", rc.handleUnmask)
	}

	forms := rc.router.Group("/forms")
	{
		forms.GET("/tokenise", rc.handleTokeniseForm)
		forms.GET("/detokenise", rc.handleDetokeniseForm)
		forms.GET("/dashboard", rc.handleDashboard)
		forms.POST("/unmask", rc.handleUnmask)
	}
}

func (rc *routerConfig) StartAPI() {
	port := os.Getenv("PORT")
	rc.router.LoadHTMLGlob("router/templates/*")
	rc.router.Static("/static", "router/static")
	rc.setupRoutes()

	for _, route := range rc.router.Routes() {
		log.Info().Msg(fmt.Sprintf("%-4s\t%s", route.Method, route.Path))
	}

	log.Info().Msgf("Attempting run at port: %s", port)

	if err := rc.router.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server")
	}
}

func (rc *routerConfig) handlePing(c *gin.Context) {
	c.String(200, "pong")
}

func (rc *routerConfig) handleTokenise(c *gin.Context) {
	ctx := context.Background()

	var payload models.TokenisePayload
	if err := c.BindJSON(&payload); err != nil {
		log.Error().Err(err).Msg("Error parsing payload for tokenise")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload."})
		return
	}

	var err error
	payload.Card.CardNumberEncrypted, err = rc.cipher.Encrypt([]byte(payload.Card.CardNumber))
	if err != nil {
		log.Error().Err(err).Str("requestID", payload.RequestID).Msg("Error encrypting credit card number")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encryption failure"})
		return
	}

	payload.Card.ExpiryDateEncrypted, err = rc.cipher.Encrypt([]byte(payload.Card.ExpiryDate))
	if err != nil {
		log.Error().Err(err).Str("requestID", payload.RequestID).Msg("Error encrypting expiry date")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encryption failure"})
		return
	}

	response, err := rc.handler.HandleTokenise(ctx, payload)
	if err != nil {
		log.Error().Err(err).Str("requestID", payload.RequestID).Msg("Error during tokenisation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tokenisation failure"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (rc *routerConfig) handleDetokenise(c *gin.Context) {
	ctx := context.Background()

	var payload models.DetokenisePayload
	if err := c.BindJSON(&payload); err != nil {
		log.Error().Err(err).Msg("Error parsing payload for detokenise")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload."})
		return
	}

	response, err := rc.handler.HandleDetokenise(ctx, payload)
	if err != nil {
		log.Error().Err(err).Str("token", payload.Token.String()).Msg("Error during detokenisation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if len(response.Card.CardNumberEncrypted) > 0 {
		cardNumberByte, err := rc.cipher.Decrypt(response.Card.CardNumberEncrypted)
		if err != nil {
			log.Error().Err(err).Str("token", payload.Token.String()).Msg("Error decrypting card number")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failure"})
			return
		}
		response.Card.CardNumber = string(cardNumberByte)
	}

	if len(response.Card.ExpiryDateEncrypted) > 0 {
		expiryDateByte, err := rc.cipher.Decrypt(response.Card.ExpiryDateEncrypted)
		if err != nil {
			log.Error().Err(err).Str("token", payload.Token.String()).Msg("Error decrypting expiry date")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failure"})
			return
		}
		response.Card.ExpiryDate = string(expiryDateByte)
	}
	if len(response.Card.ExpiryDate) > 0 {
		response.Card.ExpiryDate = response.Card.ExpiryDate[:2] + "/" + response.Card.ExpiryDate[2:]
	}

	c.JSON(http.StatusOK, response)
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
			data[i].CardNumber = common.MaskLeft(string(cardNumberDecryptedByte), 4)
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

	cards, err := rc.handler.GetAllCards(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve all cards")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cards"})
		return
	}

	if payload.Mask {
		for i, card := range cards {
			cardNumberDecryptedByte, err := rc.cipher.Decrypt(card.CardNumberEncrypted)
			if err != nil {
				log.Error().Err(err).Str("cardID", card.Token.String()).Msg("Failed to decrypt card number")
				continue
			}
			expiryDateDecryptedByte, err := rc.cipher.Decrypt(card.ExpiryDateEncrypted)
			if err != nil {
				log.Error().Err(err).Str("cardID", card.Token.String()).Msg("Failed to decrypt expiry date")
				continue
			}
			cards[i].CardNumber = common.MaskLeft(string(cardNumberDecryptedByte), 4)
			expiryDate := string(expiryDateDecryptedByte)
			cards[i].ExpiryDate = expiryDate[:2] + "/" + expiryDate[2:]
			cards[i].CardNumberEncrypted = nil
			cards[i].ExpiryDateEncrypted = nil
		}
	} else {
		for i, card := range cards {
			log.Info().Msg("Decrypting Card Number")
			cardNumberDecryptedByte, err := rc.cipher.Decrypt(card.CardNumberEncrypted)
			if err != nil {
				log.Error().Err(err).Str("cardID", card.Token.String()).Msg("Failed to decrypt card number")
				continue
			}

			log.Info().Msg("Decrypting Expiry Date")
			expiryDateDecryptedByte, err := rc.cipher.Decrypt(card.ExpiryDateEncrypted)
			if err != nil {
				log.Error().Err(err).Str("cardID", card.Token.String()).Msg("Failed to decrypt expiry date")
				continue
			}

			cards[i].CardNumber = string(cardNumberDecryptedByte)
			expiryDate := string(expiryDateDecryptedByte)
			cards[i].ExpiryDate = expiryDate[:2] + "/" + expiryDate[2:]
			cards[i].CardNumberEncrypted = nil
			cards[i].ExpiryDateEncrypted = nil
		}
	}

	c.JSON(http.StatusOK, cards)
}

func (rc *routerConfig) handleUnmask(c *gin.Context) {
	ctx := context.Background()

	var payload struct {
		Token uuid.UUID `json:"token" binding:"required"`
	}

	if err := c.BindJSON(&payload); err != nil {
		log.Error().Err(err).Msg("Invalid token provided for unmasking")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	response, err := rc.handler.HandleDetokenise(ctx, models.DetokenisePayload{Token: payload.Token})
	if err != nil {
		log.Error().Err(err).Str("token", payload.Token.String()).Msg("Error during unmasking")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unmasking failure"})
		return
	}

	if len(response.Card.CardNumberEncrypted) > 0 {
		cardNumberByte, err := rc.cipher.Decrypt(response.Card.CardNumberEncrypted)
		if err != nil {
			log.Error().Err(err).Str("token", payload.Token.String()).Msg("Error decrypting card number")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failure"})
			return
		}
		response.Card.CardNumber = string(cardNumberByte)
	}

	if len(response.Card.ExpiryDateEncrypted) > 0 {
		expiryDateByte, err := rc.cipher.Decrypt(response.Card.ExpiryDateEncrypted)
		if err != nil {
			log.Error().Err(err).Str("token", payload.Token.String()).Msg("Error decrypting expiry date")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Decryption failure"})
			return
		}
		response.Card.ExpiryDate = string(expiryDateByte)
	}
	if len(response.Card.ExpiryDate) > 0 {
		response.Card.ExpiryDate = response.Card.ExpiryDate[:2] + "/" + response.Card.ExpiryDate[2:]
	}
	c.JSON(http.StatusFound, response)
}
