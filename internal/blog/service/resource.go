package service

import (
	"errors"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
	sDto "wan_go/pkg/common/dto"
	"wan_go/sdk/service"
)

type Resource struct {
	service.Service
}

func (s *Resource) Update(d *dto.SaveResourceReq) error {

	var model blog.Resource

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

func (s *Resource) Insert(d *dto.SaveResourceReq) error {
	var model blog.Resource

	data := blog.Resource{}
	d.CopyTo(&data)
	err := s.Orm.Debug().Model(&model).Create(&data).Error
	if err != nil {
		s.Log.Errorf("InsertReq error:%s", err)
		return err
	}
	return nil
}

func (s *Resource) Delete(d *dto.DelResourceReq) error {
	var err error
	tx := s.Orm.Debug().Where("path=?", d.GetPath()).
		Delete(&blog.Resource{})

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

func (s *Resource) GetResourceInfo(resources *[]blog.Resource) error {
	tx := s.Orm.Debug().Select("id, path").
		Where("path like ? and size is null", "%"+config.Config.Qiniu.Url+"%").
		Find(&resources)
	if err := tx.Error; err != nil {
		s.Log.Errorf("GetResourceInfo error: %s", err)
		return err
	}
	return nil
}

func (s *Resource) BatchUpdate(resources *[]blog.Resource) error {
	if len(*resources) > 0 {
		tx := db.Mysql().Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		var err error
		for _, res := range *resources {
			if err = tx.Updates(res).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		return tx.Commit().Error
	}
	return nil
}

func (s *Resource) GetImageList(resources *[]blog.Resource, adminId int) error {
	tx := s.Orm.Debug().Select("path").Where("type=? and status=? and user_id=?",
		blog_const.PATH_TYPE_INTERNET_MEME, blog_const.STATUS_ENABLE.Code, adminId).
		Order("created_at DESC").Find(&resources)
	if err := tx.Error; err != nil {
		s.Log.Errorf("GetImageList error: %s", err)
		return err
	}
	return nil
}

func (s *Resource) Page(d *dto.PageResourceReq, p *actions.DataPermission, page *vo.Page[blog.Resource]) error {

	var data blog.Resource
	err := s.Orm.Debug().Scopes(
		sDto.MakeCondition(d.GetNeedSearch()),
		sDto.Paginate(d.Pagination),
		actions.Permission(data.TableName(), p),
	).Find(&page.Records).
		Limit(-1).
		Offset(-1).
		Count(&page.Total).Error

	if err != nil {
		s.Log.Errorf("Page error: %s", err)
		return err
	}

	return nil
}

func (s *Resource) ChangeResourceStatus(id int32, status bool) error {
	return s.Orm.Model(&blog.Resource{ID: id}).Update("status", status).Error
}
