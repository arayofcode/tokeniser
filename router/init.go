package router

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/arayofcode/tokeniser/handler"
)

type routerConfig struct {
	router  *gin.Engine
	handler handler.Handler
	ctx     context.Context
}

type Router interface {
	StartAPI()
}

func NewRouter(ctx context.Context, handler handler.Handler) Router {
	router := gin.Default()
	return &routerConfig{
		handler: handler,
		router:  router,
		ctx:     ctx,
	}
}
