package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stasshander/ddd/internal/infrastructure/metrics"
)

// MetricsMiddleware collects HTTP metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get status code
		status := strconv.Itoa(c.Writer.Status())

		// Record metrics
		metrics.HTTPRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
			status,
		).Inc()

		metrics.HTTPRequestDuration.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
		).Observe(duration)
	}
}
