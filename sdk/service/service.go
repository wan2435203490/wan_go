package service

import (
	"wan_go/pkg/common/api"
)

type IService interface {
	// Get 定义增删改查接口？
	Get() *Service
}

type Service struct {
	//Orm *gorm.DB
	Api *api.Api
}

func (s *Service) Get() *Service {
	return s
}

// AddError 报错即返回
func (s *Service) AddError(err error) {
	s.Api.AddError(err)
}

// AddError 报错即返回
//func (s *Service) IsError(err error) error {
//	return s.Api.IsError(err)
//}

//func (s *Service) Error(err error) {
//	r.ErrorInternal(err.Error(), s.Context)
//}
//
//func (s *Service) ErrorMsg(msg string) {
//	r.ErrorInternal(msg, s.Context)
//}
