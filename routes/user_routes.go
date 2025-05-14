package routes

import (
	"github.com/gin-gonic/gin"
	authAdapter "github.com/wittawat/go-hex/adapter/auth"
	middleware "github.com/wittawat/go-hex/adapter/middleware"
	userAdapter "github.com/wittawat/go-hex/adapter/user"
)

func RegisterUserHandler(app *gin.Engine, userHandler *userAdapter.HttpUserHandler, authNSvc *authAdapter.AuthNServiceImpl, authZSvc *authAdapter.AuthZServiceImpl) {
	public := app.Group("/")
	public.POST("/register/user", userHandler.Register)
	public.POST("/login", userHandler.Login)
	public.POST("/register/seller", userHandler.SellerRegister)

	protected := app.Group("/users", middleware.JWTAuthMiddleware(authNSvc))
	protected.GET("/:id", userHandler.GetUser)
	protected.GET("/", middleware.RequireRoles(authNSvc, authZSvc, "admin"), userHandler.GetAllUser)
	protected.PATCH("/:id", middleware.RequireRoles(authNSvc, authZSvc, "admin", "seller", "user"), userHandler.UpdateUser)
	protected.DELETE("/:id", middleware.RequireRoles(authNSvc, authZSvc, "admin", "seller", "user"), userHandler.DeleteUser)
}
