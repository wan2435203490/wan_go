package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"landlord/core/logger"
	"landlord/sdk/pkg"
	"strings"
)

// RequestId traffic
func RequestId(trafficKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetHeader(trafficKey)
		if requestId == "" {
			requestId = c.GetHeader(strings.ToLower(trafficKey))
		}
		if requestId == "" {
			requestId = uuid.New().String()
		}
		c.Request.Header.Set(trafficKey, requestId)
		c.Set(trafficKey, requestId)
		c.Set(pkg.LoggerKey,
			logger.NewHelper(logger.DefaultLogger).
				WithFields(map[string]interface{}{
					trafficKey: requestId,
				}))
		c.Next()
	}
}
