package routes

import (
	"github.com/gin-gonic/gin"
	auth "github.com/wittawat/go-hex/adapter/auth"
	middleware "github.com/wittawat/go-hex/adapter/middleware"
	adapter "github.com/wittawat/go-hex/adapter/product"
)

func RegisterProductHandler(app *gin.Engine, productHandler *adapter.HttpProductHandler, authNSvc *auth.AuthNServiceImpl, authZSvc *auth.AuthZServiceImpl) {
	protected := app.Group("/products")
	protected.GET("/", productHandler.GetAllProduct)
	protected.GET("/:id", productHandler.GetProduct)

	protected.POST("/", middleware.RequireRoles(authNSvc, authZSvc, "seller"), productHandler.CreateProduct)
	protected.PATCH("/:id", middleware.RequireRoles(authNSvc, authZSvc, "seller"), productHandler.UpdateProduct)
	protected.DELETE("/:id", middleware.RequireRoles(authNSvc, authZSvc, "sellter"), productHandler.DeleteProduct)
}
