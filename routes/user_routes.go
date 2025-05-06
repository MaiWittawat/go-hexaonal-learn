package routes

import (
	"github.com/gin-gonic/gin"
	adapter "github.com/wittawat/go-hex/adapter/user"
)

func RegisterUserRoutes(app *gin.Engine, userHandler *adapter.HttpUserHandler) {
	userRoute := app.Group("users")
	userRoute.POST("/", userHandler.Register)
	userRoute.GET("/", userHandler.GetAllUser)
	userRoute.GET("/:id", userHandler.GetUser)
	userRoute.PATCH("/:id", userHandler.UpdateUser)
	userRoute.DELETE("/:id", userHandler.DeleteUser)
}
