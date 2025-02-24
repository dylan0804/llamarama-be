package middleware

import (
	"net/http"
	"strings"

	"github.com/dylan0804/Llamarama/cmd/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(sessionStore *utils.SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}

		actualToken := strings.TrimPrefix(token, "Bearer ")

		user, err := sessionStore.Get(c.Request.Context(), actualToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user_id", user["user_id"])

		c.Next()
	}
}