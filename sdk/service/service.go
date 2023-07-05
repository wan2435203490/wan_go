package service

import (
	"fmt"
	"gorm.io/gorm"
	"wan_go/core/logger"
)

type IService interface {
	Get() *Service
}

func (s *Service) AddError(err error) error {
	if s.Error == nil {
		s.Error = err
	} else if err != nil {
		s.Error = fmt.Errorf("%v; %w", s.Error, err)
	}
	return s.Error
}

type Service struct {
	Orm   *gorm.DB
	Msg   string
	MsgID string
	Log   *logger.Helper
	Error error
}

func (s *Service) Get() *Service {
	return s
}
