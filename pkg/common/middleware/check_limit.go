package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"wan_go/sdk/pkg/response"
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

		response.Error(c, http.StatusInternalServerError, nil, "请求过于频繁")
		//res := &r.Response{}
		//res.Message = "请求过于频繁"
		//res.Status = r.ErrorStatus
		//res.Captcha = http.StatusInternalServerError
		//c.AbortWithStatusJSON(http.StatusTooManyRequests, res)
		return
	}
}
