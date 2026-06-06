package middleware

import (
	"net/http"
	"strings"

	"Distribyte/backend/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader :=
			c.GetHeader("Authorization")

		if authHeader == "" {

			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Missing token",
			})

			c.Abort()
			return
		}

		tokenString :=
			strings.TrimPrefix(
				authHeader,
				"Bearer ",
			)

		claims, err :=
			utils.ValidateToken(
				tokenString,
			)

		if err != nil {

			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token",
			})

			c.Abort()
			return
		}

		c.Set(
			"user_id",
			claims.UserID,
		)

		c.Set(
			"email",
			claims.Email,
		)

		c.Next()
	}
}
