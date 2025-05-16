package middleware

import (
	"eshop/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пытаемся получить токен из cookie
		tokenString, err := c.Cookie("Authorization")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		userID, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

// func JWTAuth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
// 			return
// 		}

// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		userID, err := utils.ParseJWT(tokenString)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
// 			return
// 		}

// 		// сохраняем userID в контексте
// 		c.Set("userID", userID)
// 		c.Next()
// 	}
// }

// func JWTAuth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var tokenString string

// 		// 1) сначала смотрим заголовок Authorization
// 		authHeader := c.GetHeader("Authorization")
// 		if strings.HasPrefix(authHeader, "Bearer ") {
// 			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
// 		} else {
// 			// 2) если нет — пытаемся взять cookie "auth"
// 			if cookie, err := c.Cookie("auth"); err == nil {
// 				tokenString = cookie
// 			}
// 		}

// 		if tokenString == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized,
// 				gin.H{"error": "missing or invalid Authorization header"})
// 			return
// 		}

// 		userID, err := utils.ParseJWT(tokenString)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized,
// 				gin.H{"error": "invalid token"})
// 			return
// 		}

// 		c.Set("userID", userID) // сохраняем id в контексте
// 		c.Next()
// 	}
// }
