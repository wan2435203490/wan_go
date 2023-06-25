package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	r "wan_go/pkg/common/response"
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
		res := &r.Response{}
		res.Message = "请求过于频繁"
		res.Status = r.ErrorStatus
		res.Code = http.StatusInternalServerError
		c.AbortWithStatusJSON(http.StatusTooManyRequests, res)
		return
	}
}
