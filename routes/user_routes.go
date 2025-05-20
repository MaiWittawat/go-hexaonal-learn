package routes

import (
	"github.com/gin-gonic/gin"
	authNAdapter "github.com/wittawat/go-hex/adapter/auth"
	middleware "github.com/wittawat/go-hex/adapter/middleware"
	userAdapterInbound "github.com/wittawat/go-hex/adapter/user/inbound"
	authZSvc "github.com/wittawat/go-hex/core/services"
)

func RegisterUserHandler(app *gin.Engine, userHandler *userAdapterInbound.HttpUserHandler, authNAdapter *authNAdapter.AuthenService, authZSvc *authZSvc.AuthorizeService) {
	public := app.Group("/")
	public.POST("/register/user", userHandler.Register)
	public.POST("/login", userHandler.Login)
	public.POST("/register/seller", userHandler.SellerRegister)

	protected := app.Group("/users", middleware.AuthenticationMiddleware(authNAdapter))
	protected.GET("/:id", userHandler.GetUser)
	protected.GET("/", middleware.AuthorizeRoles(authNAdapter, authZSvc, "admin"), userHandler.GetAllUser)
	protected.PATCH("/:id", middleware.AuthorizeRoles(authNAdapter, authZSvc, "admin", "seller", "user"), userHandler.UpdateUser)
	protected.DELETE("/:id", middleware.AuthorizeRoles(authNAdapter, authZSvc, "admin", "seller", "user"), userHandler.DeleteUser)
}
