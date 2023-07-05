package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"net/http"
	"wan_go/core/logger"
	"wan_go/sdk/pkg"
	"wan_go/sdk/service"
)

const (
	DefaultLanguage = "zh-CN"
)

type Api struct {
	Context *gin.Context
	Logger  *logger.Helper
	Orm     *gorm.DB
	Errs    error
}

func (a *Api) AddError(err error) {
	if err == nil {
		return
	}
	if a.Errs == nil {
		a.Errs = err
	} else {
		a.Logger.Error(err)
		a.Errs = fmt.Errorf("%v; %w", a.Errs, err)
	}
}

// IsError 如果有error则write json response
func (a *Api) IsError(err error) bool {
	if err != nil {
		a.ErrorInternal(err.Error())
		return true
	}
	return false
}

//func (a Api) Build(c *gin.Context, s service.IService, d interface{}, bindings ...binding.Binding) error {
//	return a.MakeContext(c).MakeOrm().Bind(d, bindings...).MakeService(s.Get()).CodeMsg
//}

func (a *Api) MakeContextChain(c *gin.Context, s *service.Service, d interface{}, bindings ...binding.Binding) *Api {
	err := a.MakeContext(c).MakeOrm().Binds(d, bindings...).MakeService(s).Errs
	if err != nil {
		a.Logger.Error(err)
		a.ErrorInternal(err.Error())
		return nil
	}
	return a
}

// MakeContext 设置http上下文
func (a *Api) MakeContext(c *gin.Context) *Api {
	a.Context = c
	a.Logger = GetRequestLogger(c)
	return a
}

// GetLogger 获取上下文提供的日志
func (a Api) GetLogger() *logger.Helper {
	return GetRequestLogger(a.Context)
}

// GetOrm 获取Orm DB
func (a Api) GetOrm() (*gorm.DB, error) {
	db, err := pkg.GetOrm(a.Context)
	if err != nil {
		a.Logger.Error(http.StatusInternalServerError, err, "数据库连接获取失败")
		return nil, err
	}
	return db, nil
}

// MakeOrm 设置Orm DB
func (a *Api) MakeOrm() *Api {
	var err error
	if a.Logger == nil {
		err = errors.New("at MakeOrm logger is nil")
		a.AddError(err)
		return a
	}
	db, err := pkg.GetOrm(a.Context)
	if err != nil {
		a.Logger.Error(http.StatusInternalServerError, err, "数据库连接获取失败")
		a.AddError(err)
	}
	a.Orm = db
	return a
}

func (a *Api) MakeService(s *service.Service) *Api {
	if s == nil {
		return a
	}
	s.Log = a.Logger
	s.Orm = a.Orm
	return a
}

func (a Api) GetRequest() *http.Request {
	return a.Context.Request
}
