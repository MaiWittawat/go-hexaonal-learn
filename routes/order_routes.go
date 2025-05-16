package routes

import (
	"github.com/gin-gonic/gin"
	authNAdapter "github.com/wittawat/go-hex/adapter/auth"
	middlewareAdapter "github.com/wittawat/go-hex/adapter/middleware"
	orderAdapterInbound "github.com/wittawat/go-hex/adapter/order/inbound"
	authZSvc "github.com/wittawat/go-hex/core/service"
)

func RegisterOrderHandler(app *gin.Engine, orderHandler *orderAdapterInbound.HttpOrderHandler, authNAdapter *authNAdapter.AuthenService, authZSvc *authZSvc.AuthorizeService) {
	protected := app.Group("/orders", middlewareAdapter.JWTAuthMiddleware(authNAdapter))
	protected.GET("/user/:user_id", middlewareAdapter.RequireRoles(authNAdapter, authZSvc, "user", "seller", "admin"), orderHandler.FindOrder)
	protected.POST("/", middlewareAdapter.RequireRoles(authNAdapter, authZSvc, "user", "seller", "admin"), orderHandler.CreateOrder)
	protected.PATCH("/:id", middlewareAdapter.RequireRoles(authNAdapter, authZSvc, "user", "seller", "admin"), orderHandler.UpdateOrder)
	protected.DELETE("/:id", middlewareAdapter.RequireRoles(authNAdapter, authZSvc, "user", "seller", "admin"), orderHandler.DeleteOrder)
}
