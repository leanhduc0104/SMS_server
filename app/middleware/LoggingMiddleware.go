package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(queryLog *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy thời gian bắt đầu
		startTime := time.Now()

		// Xử lý truy vấn
		c.Next()

		// Lấy thời gian kết thúc
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// Lấy thông tin truy vấn
		method := c.Request.Method
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		path := c.Request.URL.Path

		// Ghi log thông tin truy vấn
		queryLog.Printf("| %3d | %13v | %15s | %-7s  %#v\n",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
