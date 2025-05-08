package routes

import (
	"github.com/gin-gonic/gin"
	auth "github.com/wittawat/go-hex/adapter/auth"
	middleware "github.com/wittawat/go-hex/adapter/middleware"
	adapter "github.com/wittawat/go-hex/adapter/user"
)

func RegisterUserRoutes(app *gin.Engine, userHandler *adapter.HttpUserHandler, jwtSvc *auth.AuthNServiceImpl) {
	public := app.Group("/")
	public.POST("/register", userHandler.Register)
	public.POST("/login", userHandler.Login)

	protected := app.Group("/users", middleware.JWTAuthMiddleware(jwtSvc))
	protected.GET("/", userHandler.GetAllUser)
	protected.GET("/:id", userHandler.GetUser)
	protected.PUT("/:id", userHandler.UpdateUser)
	protected.DELETE("/:id", userHandler.DeleteUser)
}
