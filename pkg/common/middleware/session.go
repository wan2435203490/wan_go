package middleware

import (
	json2 "encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"landlord/common/config"
	"landlord/common/token"
	"landlord/db"
	r "landlord/db_common/response"
	"strings"
)

func WithSession(c *gin.Context) {

	if strings.HasPrefix(c.Request.RequestURI, "/auth") {
		//strings.HasPrefix(c.Request.RequestURI, "/get") ||
		//strings.HasPrefix(c.Request.RequestURI, "/set")
		c.Next()
		return
	}

	var user db.User
	session := sessions.Default(c)
	get := session.Get(config.Config.Session.UserSessionKey)

	if get == nil {
		//session没有的话 用jwt解析Header.Token
		tokenStr := c.GetHeader("Token")
		ok, userId := token.GetUserIdFromToken(tokenStr)
		if !ok {
			c.Abort()
			r.Error(401, "session记录被清除，请退出重新登录", c)
			return
		}
		user.Id = userId
		if db.MySQL.Find(&user).Error != nil {
			c.Abort()
			r.Error(401, "用户不存在，请重新登录", c)
			return
		}
		c.Set(config.Config.Session.UserSessionKey, &user)
		c.Next()
		return
	}

	err := json2.Unmarshal(get.([]byte), &user)
	if err != nil {
		c.Abort()
		r.Error(401, "session记录被清除，请退出重新登录", c)
		return
	}

	c.Set(config.Config.Session.UserSessionKey, &user)
	c.Next()
}
