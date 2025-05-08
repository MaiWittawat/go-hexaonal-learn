package routes

import (
	"github.com/gin-gonic/gin"
	authAdapter "github.com/wittawat/go-hex/adapter/auth"
	middleware "github.com/wittawat/go-hex/adapter/middleware"
	userAdapter "github.com/wittawat/go-hex/adapter/user"
)

func RegisterUserRoutes(app *gin.Engine, userHandler *userAdapter.HttpUserHandler, authNSvc *authAdapter.AuthNServiceImpl, authZSvc *authAdapter.AuthZServiceImpl) {
	public := app.Group("/")
	public.POST("/register", userHandler.Register)
	public.POST("/login", userHandler.Login)

	protected := app.Group("/users", middleware.JWTAuthMiddleware(authNSvc))
	protected.GET("/", middleware.RequireRoles(authNSvc, authZSvc, "admin"), userHandler.GetAllUser)
	protected.GET("/:id", userHandler.GetUser)
	protected.PATCH("/:id", userHandler.UpdateUser)
	protected.DELETE("/:id", userHandler.DeleteUser)
}
