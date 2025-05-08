package routes

import (
	"github.com/gin-gonic/gin"
	authAdapter "github.com/wittawat/go-hex/adapter/auth"
	middlewareAdapter "github.com/wittawat/go-hex/adapter/middleware"
	orderAdapter "github.com/wittawat/go-hex/adapter/order"
)

func RegisterOrderHandler(app *gin.Engine, orderHandler *orderAdapter.HttpOrderHandler, authNSvc *authAdapter.AuthNServiceImpl, authZSvc *authAdapter.AuthZServiceImpl) {
	protected := app.Group("/orders", middlewareAdapter.JWTAuthMiddleware(authNSvc))
	protected.GET("/user/:user_id", orderHandler.FindOrder)
	protected.POST("/", orderHandler.CreateOrder)
	protected.PATCH("/:id", orderHandler.UpdateOrder)
	protected.DELETE("/:id", orderHandler.DeleteOrder)
}
