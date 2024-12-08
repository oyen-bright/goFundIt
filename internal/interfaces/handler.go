package interfaces

import "github.com/gin-gonic/gin"

type HandlerInterface interface {
	RegisterRoutes(router *gin.RouterGroup, middlewares []gin.HandlerFunc)
}
