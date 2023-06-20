package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"wan_go/core/logger"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/utils"
	"wan_go/sdk/pkg"
)

const (
	DefaultLanguage = "zh-CN"
)

type Api struct {
	Context *gin.Context
	//Logger  *logger.Helper
	Orm *gorm.DB
	Err error
}

func (a *Api) AddError(err error) {
	if err == nil {
		return
	}
	if a.Err == nil {
		a.Err = err
	} else {
		a.Err = fmt.Errorf("%v; %w", a.Err, err)
	}
	return
}

// IsError 如果有error则write json response
func (a *Api) IsError(err error) bool {
	if err != nil {
		a.ErrorInternal(err.Error())
		return true
	}
	return false
}

//func (a *Api) Build(c *gin.Context, s service.IService, d interface{}, bindings ...binding.Binding) error {
//	return a.MakeContext(c).MakeOrm().BindFailed(d, bindings...).MakeService(s.Get()).CodeMsg
//}

// MakeContext 设置http上下文
func (a *Api) MakeContext(c *gin.Context) error {
	a.Context = c
	//a.Logger = GetRequestLogger(c)

	//err := a.MakeOrm()
	//if err != nil {
	//	return err
	//}

	return nil
}

// MakeOrm 设置Orm DB
func (a *Api) MakeOrm() error {
	var err error
	db, err := pkg.GetOrm(a.Context)
	if err != nil {
		//s.Api.Logger.Error(http.StatusInternalServerError, err, "数据库连接获取失败")
		return err
	}
	a.Orm = db
	return nil
}

// GetLogger 获取上下文提供的日志
func (a *Api) GetLogger() *logger.Helper {
	return GetRequestLogger(a.Context)
}

//
//func (a *Api) MakeService(c *service.Service) *Api {
//	c.Log = a.Logger
//	c.Orm = a.Orm
//	c.Context = a.Context
//	return a
//}

func (a *Api) GetRequest() *http.Request {
	return a.Context.Request
}

func (a *Api) GetToken() string {
	return a.GetRequest().Header[blog_const.TOKEN_HEADER][0]
}

func (a *Api) GetCurrentUser() *blog.User {
	if get, b := cache.Get(a.GetToken()); b {
		return get.(*blog.User)
	}
	return nil
}

func (a *Api) GetCurrentUserId() int32 {
	if user := a.GetCurrentUser(); user != nil {
		return user.ID
	}

	return -1
}

func (a *Api) GetCurrentUserIdStr() string {
	if user := a.GetCurrentUser(); user != nil {
		return utils.Int32ToString(user.ID)
	}

	return ""
}

// todo user cache
func (a *Api) KeyUserId(pre string) string {
	if user := a.GetCurrentUser(); user != nil {
		return pre + utils.Int32ToString(user.ID)
	}

	return pre
}

func (a *Api) AdminId() int32 {
	admin := cache.GetAdminUser()
	if admin == nil {
		return 0
	}

	return admin.ID
}

func (a *Api) IsAdmin() bool {

	if a.AdminId() == 0 {
		return false
	}

	return a.GetCurrentUserId() == a.AdminId()
}
