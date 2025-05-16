package routes

import (
	"github.com/gin-gonic/gin"
	authNAdapter "github.com/wittawat/go-hex/adapter/auth"
	middleware "github.com/wittawat/go-hex/adapter/middleware"
	productAdapterInbound "github.com/wittawat/go-hex/adapter/product/inbound"
	authZSvc "github.com/wittawat/go-hex/core/service"
)

func RegisterProductHandler(app *gin.Engine, productHandler *productAdapterInbound.HttpProductHandler, authNAdapter *authNAdapter.AuthenService, authZSvc *authZSvc.AuthorizeService) {
	public := app.Group("/products")
	public.GET("/", productHandler.GetAllProduct)
	public.GET("/:id", productHandler.GetProduct)

	protected := app.Group("/products", middleware.JWTAuthMiddleware(authNAdapter))
	protected.POST("/", middleware.RequireRoles(authNAdapter, authZSvc, "seller", "admin"), productHandler.CreateProduct)
	protected.PATCH("/:id", middleware.RequireRoles(authNAdapter, authZSvc, "seller", "admin"), productHandler.UpdateProduct)
	protected.DELETE("/:id", middleware.RequireRoles(authNAdapter, authZSvc, "seller", "admin"), productHandler.DeleteProduct)
}
