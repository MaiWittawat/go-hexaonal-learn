package routes

import (
	"github.com/gin-gonic/gin"
	adapter "github.com/wittawat/go-hex/adapter/order"
)

func RegisterOrderHandler(app *gin.Engine, orderHandler *adapter.HttpOrderHandler) {
	orderRoute := app.Group("orders")
	orderRoute.GET("/user/:user_id", orderHandler.FindOrder)
	orderRoute.POST("/", orderHandler.CreateOrder)
	orderRoute.PATCH("/:id", orderHandler.UpdateOrder)
	orderRoute.DELETE("/:id", orderHandler.DeleteOrder)
}
