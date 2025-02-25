package middleware

import (
	"net/http"
	"strings"

	"github.com/dylan0804/Llamarama/cmd/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(sessionStore *utils.SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		authHeader := c.GetHeader("Authorization")

		if authHeader != "" {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			token = c.Query("token")
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}

		user, err := sessionStore.Get(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user_id", user["user_id"])

		c.Next()
	}
}