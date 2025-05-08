// adapter/middleware/jwt_authZ.go
package adapter

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	port "github.com/wittawat/go-hex/core/port/auth"
)

func RequireRoles(tokenSvc port.JwtAuthNService, authZSvc port.JwtAuthZService, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			log.Println("missing token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claimsVal, exists := c.Get("claims")
		if !exists {
			log.Println("missing claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims := claimsVal.(jwt.MapClaims)
		email, ok := claims["email"].(string)
		if !ok {
			log.Println("invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		authorized, err := authZSvc.Authorize(email, roles)
		if !authorized || err != nil {
			log.Println("invalid token claims: ", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("userEmail", email)
		c.Next()
	}
}
