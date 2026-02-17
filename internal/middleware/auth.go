// Package middleware contains HTTP middleware implementations.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tahiriqbal095/attendify/internal/models"
	"github.com/tahiriqbal095/attendify/internal/service"
)

const (
	ContextKeyUserID = "user_id"
	ContextKeyRole   = "role"
)

func Auth(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			abortUnauthorized(c, "authorization header required")
			return
		}

		// Expect format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			abortUnauthorized(c, "invalid authorization header format")
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			abortUnauthorized(c, "token required")
			return
		}

		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			abortUnauthorized(c, "invalid or expired token")
			return
		}

		// Store user info in context for downstream handlers
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyRole, claims.Role)

		c.Next()
	}
}

func GetUserID(c *gin.Context) uuid.UUID {
	if id, exists := c.Get(ContextKeyUserID); exists {
		if userID, ok := id.(uuid.UUID); ok {
			return userID
		}
	}
	return uuid.Nil
}

func GetUserRole(c *gin.Context) models.Role {
	if role, exists := c.Get(ContextKeyRole); exists {
		if r, ok := role.(models.Role); ok {
			return r
		}
	}
	return ""
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error":   message,
	})
}
