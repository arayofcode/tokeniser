package router

import (
	"net/http"

	"github.com/arayofcode/tokeniser/models"
	"github.com/gin-gonic/gin"
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
	var payload models.TokenisePayload
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload."})
	}

	response := rc.handler.HandleTokeniseNew(rc.ctx, payload)
	c.JSON(http.StatusCreated, response)
}

func (rc *routerConfig) handleDetokenise(c *gin.Context) {
	var payload models.DetokenisePayload
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload."})
	}
	response := rc.handler.HandleDetokeniseNew(rc.ctx, payload)
	c.JSON(http.StatusFound, response)
}
