package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stasshander/ddd/internal/infrastructure/metrics"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

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
