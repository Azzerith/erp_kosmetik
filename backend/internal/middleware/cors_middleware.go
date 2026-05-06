package middleware

import (
	"erp-cosmetics-backend/internal/config"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origins := strings.Split(cfg.CORSAllowedOrigins, ",")
		origin := c.GetHeader("Origin")

		for _, o := range origins {
			if o == "*" || o == origin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", cfg.CORSAllowedMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", cfg.CORSAllowedHeaders)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}