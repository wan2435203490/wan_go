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

type Sort struct {
	service.Service
}

func NewSort(c *gin.Context) *Sort {
	us := Sort{}
	us.Orm = pkg.Orm(c)
	us.Log = api.GetRequestLogger(c)
	return &us
}

func (s *Sort) Update(d *dto.SaveSortReq) error {

	var model blog.Sort
	var count int64
	if err := s.Orm.Debug().Model(&model).Where("id<>? and sort_name=?", d.ID, d.SortName).Count(&count).Error; err != nil {
		s.Log.Errorf("Update Count error:%s", err)
		return err
	}
	if count > 0 {
		return errors.New("分类名称不能重复！")
	}

	d.CopyTo(&model)
	tx := s.Orm.Debug().Model(&model).Updates(&d)
	if err := tx.Error; err != nil {
		s.Log.Errorf("Update error:%s", err)
		return err
	}
	if tx.RowsAffected == 0 {
		return errors.New("数据不存在或无权更新该数据")
	}
	return nil
}

func (s *Sort) Insert(d *dto.SaveSortReq) error {

	var model blog.Sort
	var count int64
	if err := s.Orm.Debug().Model(&model).Where("sort_name=?", d.SortName).Count(&count).Error; err != nil {
		s.Log.Errorf("InsertReq Count error:%s", err)
		return err
	}

	if count > 0 {
		return errors.New("分类名称不能重复！")
	}

	data := blog.Sort{}
	d.CopyTo(&data)
	err := s.Orm.Debug().Model(&model).Create(&data).Error
	if err != nil {
		s.Log.Errorf("InsertReq error:%s", err)
		return err
	}
	return nil
}

func (s *Sort) Delete(d *dto.DelSortReq) error {
	var err error
	var data blog.Sort

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

func (s *Sort) Page(d *dto.PageSortReq, p *actions.DataPermission, page *vo.Page[blog.Sort]) error {

	var data blog.Sort
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
