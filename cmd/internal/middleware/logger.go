package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		errors := c.Errors.Errors()

		latency := time.Since(start)

		attrs := []slog.Attr{
			slog.String("request_id", requestID),
            slog.String("method", c.Request.Method),
            slog.String("path", path),
            slog.Int("status", c.Writer.Status()),
            slog.Duration("latency", latency),
            slog.String("ip", c.ClientIP()),
            slog.String("user_agent", c.Request.UserAgent()),
            slog.Int("body_size", c.Writer.Size()),
		}

		if raw != "" {
			attrs = append(attrs, slog.String("query", raw))
		}

		if len(errors) > 0 {
			attrs = append(attrs, slog.Any("errors", errors))
		}
		
		if c.Writer.Status() >= 500 {
            slog.Error("server error", slog.Any("attrs", attrs))
        } else if c.Writer.Status() >= 400 {
            slog.Warn("client error", slog.Any("attrs", attrs))
        } else {
            slog.Info("request completed", slog.Any("attrs", attrs))
        }
	}
}
