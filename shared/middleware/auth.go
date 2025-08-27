package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"ua/shared/auth"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := bearerToken[1]
		claims, err := auth.ValidateToken(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		ctx := context.WithValue(c.Request.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func OptionalAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.Next()
			return
		}

		token := bearerToken[1]
		claims, err := auth.ValidateToken(token, jwtSecret)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)

			ctx := context.WithValue(c.Request.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "username", claims.Username)
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}

// WebSocketAuthMiddleware supports both Authorization header and token query parameter
func WebSocketAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// Try Authorization header first
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) == 2 && bearerToken[0] == "Bearer" {
				token = bearerToken[1]
			}
		}

		// If no Authorization header, try query parameter
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		ctx := context.WithValue(c.Request.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
