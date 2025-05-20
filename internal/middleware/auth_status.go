package middleware

import (
	"github.com/gin-gonic/gin"
)

func AddLoginStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		isLoggedIn := exists && userID != nil
		c.Set("isLoggedIn", isLoggedIn)
		c.Next()
	}
}
