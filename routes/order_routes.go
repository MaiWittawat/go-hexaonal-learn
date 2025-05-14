package routes

import (
	"github.com/gin-gonic/gin"
	authAdapter "github.com/wittawat/go-hex/adapter/auth"
	middlewareAdapter "github.com/wittawat/go-hex/adapter/middleware"
	orderAdapter "github.com/wittawat/go-hex/adapter/order"
)

func RegisterOrderHandler(app *gin.Engine, orderHandler *orderAdapter.HttpOrderHandler, authNSvc *authAdapter.AuthNServiceImpl, authZSvc *authAdapter.AuthZServiceImpl) {
	protected := app.Group("/orders", middlewareAdapter.JWTAuthMiddleware(authNSvc))
	protected.GET("/user/:user_id", middlewareAdapter.RequireRoles(authNSvc, authZSvc, "user", "seller", "admin"), orderHandler.FindOrder)
	protected.POST("/", middlewareAdapter.RequireRoles(authNSvc, authZSvc, "user", "seller", "admin"), orderHandler.CreateOrder)
	protected.PATCH("/:id", middlewareAdapter.RequireRoles(authNSvc, authZSvc, "user", "seller", "admin"), orderHandler.UpdateOrder)
	protected.DELETE("/:id", middlewareAdapter.RequireRoles(authNSvc, authZSvc, "user", "seller", "admin"), orderHandler.DeleteOrder)
}
