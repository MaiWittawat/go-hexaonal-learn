package routes

import (
	"github.com/gin-gonic/gin"
	adapter "github.com/wittawat/go-hex/adapter/product"
)

func RegisterProductHandler(app *gin.Engine, productHandler *adapter.HttpProductHandler) {
	productRote := app.Group("products")
	productRote.GET("/", productHandler.GetAllProduct)
	productRote.GET("/:id", productHandler.GetProduct)
	productRote.POST("/", productHandler.CreateProduct)
	productRote.PATCH("/:id", productHandler.UpdateProduct)
	productRote.DELETE("/:id", productHandler.DeleteProduct)
}
