package middleware

import (
	"net/http"

	"erp-cosmetics-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid role", nil)
			c.Abort()
			return
		}

		for _, role := range roles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", nil)
		c.Abort()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin", "super_admin")
}

func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole("super_admin")
}