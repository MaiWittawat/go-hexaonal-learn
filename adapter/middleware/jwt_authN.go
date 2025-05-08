package adapter

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	auth "github.com/wittawat/go-hex/core/port/auth"
)

func JWTAuthMiddleware(tokenService auth.JwtAuthNService) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.GetHeader("Authorization")
		token := strings.TrimPrefix(s, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, err := tokenService.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
