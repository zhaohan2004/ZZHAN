package middleware

import (
	"fmt"
	"time"

	"ZZHAN/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", status),
			zap.String("latency", formatLatency(latency)),
		}

		if status >= 400 {
			if len(c.Errors) > 0 {
				fields = append(fields, zap.String("error", c.Errors.String()))
			}
			logger.Error("请求失败", fields...)
		} else {
			logger.Info("请求成功", fields...)
		}
	}
}

// formatLatency 格式化耗时，如 6.6ms、1.2s
func formatLatency(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.1fms", float64(d.Microseconds())/1000)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}
