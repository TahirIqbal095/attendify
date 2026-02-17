package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tahiriqbal095/attendify/internal/models"
)

func RequireRole(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := GetUserRole(c)
		if userRole == "" {
			abortForbidden(c, "role not found in context")
			return
		}

		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		abortForbidden(c, "insufficient permissions")
	}
}

func RequireTeacher() gin.HandlerFunc {
	return RequireRole(models.RoleTeacher)
}

func RequireStudent() gin.HandlerFunc {
	return RequireRole(models.RoleStudent)
}

func abortForbidden(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"success": false,
		"error":   message,
	})
}
