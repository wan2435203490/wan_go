package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

var rateLimiterMap sync.Map

func WithLimit(c *gin.Context) {

	session := sessions.Default(c)
	var sessionId = session.ID()

	limiter, ok := rateLimiterMap.Load(sessionId)
	if !ok {
		//todo config
		limiter = rate.NewLimiter(5, 20)
		rateLimiterMap.Store(sessionId, limiter)
	}

	if !limiter.(*rate.Limiter).Allow() {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, "请求过于频繁")
	}
}
