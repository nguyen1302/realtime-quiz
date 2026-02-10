package middleware

import (
	"net/http"
	"strings"

	"github.com/nguyen1302/realtime-quiz/internal/service"
	"github.com/nguyen1302/realtime-quiz/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString = parts[1]
			}
		}

		// Fallback to query parameter (often used for WebSockets)
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization required", nil)
			c.Abort()
			return
		}

		// Validate token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			c.Abort()
			return
		}

		// Store claims in context for later use
		c.Set("claims", claims)
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
