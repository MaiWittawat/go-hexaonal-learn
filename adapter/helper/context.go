package helper

import (
	"github.com/gin-gonic/gin"
)

func GetUserEmail(c *gin.Context) (string, bool) {
	userEmailVal, exists := c.Get("userEmail")
	if !exists {
		return "", false
	}
	email, ok := userEmailVal.(string)
	return email, ok
}
