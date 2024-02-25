package router

import (
	"github.com/gin-gonic/gin"

	"github.com/arayofcode/tokeniser/handler"
)

type routerConfig struct {
	router  *gin.Engine
	handler handler.Handler
}

type Router interface {
	StartAPI()
}

func NewRouter(handler handler.Handler) Router {
	router := gin.Default()
	return &routerConfig{
		handler: handler,
		router:  router,
	}
}
