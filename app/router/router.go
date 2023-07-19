package router

import "github.com/gin-gonic/gin"

type IRouter interface {
	Register(engine *gin.Engine) error
}

type Router struct {
}
