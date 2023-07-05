package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/db/mysql/blog"
	sDto "wan_go/pkg/common/dto"
	"wan_go/sdk/pkg"
	"wan_go/sdk/service"
)

type Family struct {
	service.Service
}

func NewFamily(c *gin.Context) *Family {
	us := Family{}
	us.Orm = pkg.Orm(c)
	us.Log = api.GetRequestLogger(c)
	return &us
}

func (s *Family) GetByUserId(userId int32, family *blog.Family) error {
	var err error
	if err = s.Orm.Debug().Where("user_id=?", userId).First(&family).Error; err != nil {
		s.Log.Errorf("GetByUserId error:%s", err)
		return err
	}
	return nil
}

func (s *Family) Update(d *blog.Family) error {
	tx := s.Orm.Debug().Model(d).Updates(&d)
	if err := tx.Error; err != nil {
		s.Log.Errorf("Update error:%s", err)
		return err
	}
	if tx.RowsAffected == 0 {
		return errors.New("数据不存在或无权更新该数据")
	}
	return nil
}

func (s *Family) Insert(d *blog.Family) error {

	err := s.Orm.Debug().Model(d).Create(&d).Error
	if err != nil {
		s.Log.Errorf("InsertReq error:%s", err)
		return err
	}
	return nil
}

func (s *Family) Delete(d *dto.DelFamilyReq) error {
	var err error
	var data blog.Family

	tx := s.Orm.Debug().Model(&data).Delete(&data, d.GetId())

	if err = tx.Error; err != nil {
		err = tx.Error
		s.Log.Errorf("Delete error: %s", err)
		return err
	}
	if tx.RowsAffected == 0 {
		err = errors.New("数据不存在或无权删除该数据")
		return err
	}
	return nil
}

func (s *Family) Page(d *dto.PageFamilyReq, p *actions.DataPermission, page *vo.Page[blog.Family]) error {

	var data blog.Family
	err := s.Orm.Debug().Scopes(
		sDto.MakeCondition(d.GetNeedSearch()),
		sDto.Paginate(d.Pagination),
		actions.Permission(data.TableName(), p),
	).Find(&page.Records).
		Limit(-1).
		Offset(-1).
		Count(&page.Total).Error

	if err != nil {
		s.Log.Errorf("PageUser error: %s", err)
		return err
	}

	return nil
}

func (s *Family) ChangeLoveStatus(d *dto.ChangeFamilyReq) error {
	return s.Orm.Debug().Model(&blog.Family{ID: d.ID}).Update("status", d.Status).Error
}
