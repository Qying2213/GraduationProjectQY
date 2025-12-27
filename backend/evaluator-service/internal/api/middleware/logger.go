package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"evaluator-service/internal/logging"
)

func Logger(log *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start)
		log.Info("req", logging.KV("method", c.Request.Method), logging.KV("path", c.Request.URL.Path), logging.KV("status", c.Writer.Status()), logging.KV("latency", dur.String()), logging.KV("rid", c.GetString("request_id")))
	}
}

