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

	response, err := rc.handler.HandleTokenise(ctx, payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
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

	c.JSON(http.StatusFound, response)
}
