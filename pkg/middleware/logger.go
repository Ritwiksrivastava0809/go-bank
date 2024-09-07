package middleware

import (
	"time"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		userID := c.Request.Header.Get(constants.UserID)
		// Process the request
		c.Next()

		// Log the request and response details
		end := time.Now()
		latency := end.Sub(start)

		log.Info().
			Msgf("Handled request %s=%v %s=%s %s=%s %s=%d %s=%s",
				"latency", latency,
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", c.Writer.Status(),
				"user_id", userID,
			)

	}
}
