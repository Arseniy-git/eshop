package middleware

import (
	"eshop/internal/utils"

	"github.com/gin-gonic/gin"
)

// Эта middleware будет пытаться достать токен из cookie и установить userID в контекст
func JWTFromCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("Authorization")
		if err == nil && token != "" {
			userID, err := utils.ParseJWT(token)
			if err == nil {
				c.Set("userID", userID)
			}
		}
		c.Next() // всегда продолжаем
	}
}
