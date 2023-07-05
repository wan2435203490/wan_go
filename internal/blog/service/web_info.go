package service

import (
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/service/dto"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/sdk/pkg"
	"wan_go/sdk/service"
)

type WebInfo struct {
	service.Service
}

func NewWebInfo(c *gin.Context) *WebInfo {
	us := WebInfo{}
	us.Orm = pkg.Orm(c)
	us.Log = api.GetRequestLogger(c)
	return &us
}

func (s *WebInfo) Update(d *dto.SaveWebInfoReq) error {
	var webInfo blog.WebInfo
	d.CopyTo(&webInfo)
	return s.Orm.Debug().Updates(&webInfo).Error
}

func (s *WebInfo) List(webInfo *[]blog.WebInfo) error {
	if err := s.Orm.Debug().Model(&blog.WebInfo{}).Find(webInfo).Error; err != nil {
		return err
	}
	return nil
}
