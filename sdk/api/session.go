package api

import (
	. "encoding/json"
	"github.com/gin-contrib/sessions"
	"wan_go/common/config"
	"wan_go/pkg/common/db"
)

func (a *Api) User() *db.User {
	session := sessions.Default(a.Context)
	value := session.Get(config.Config.Session.UserSessionKey)

	if value != nil {
		var user db.User
		err := Unmarshal(value.([]byte), &user)
		if err != nil {
			panic(err)
		}

		return &user
	}

	user, ok := a.Context.Get(config.Config.Session.UserSessionKey)
	if !ok {
		panic("用户记录被清除，请退出重新登录")
	}
	return user.(*db.User)
}

// SetUser todo 根据type反射set any type
func (a *Api) SetUser(user *db.User) {
	session := sessions.Default(a.Context)

	userJson, err := Marshal(user)
	if err != nil {
		panic(err)
	}
	session.Set(config.Config.Session.UserSessionKey, userJson)

	_ = session.Save()
}
