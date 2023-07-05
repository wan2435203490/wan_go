package service

import (
	"errors"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/db/mysql/blog"
	sDto "wan_go/pkg/common/dto"
	"wan_go/sdk/service"
)

type Label struct {
	service.Service
}

func (s *Label) Update(d *dto.SaveLabelReq) error {

	var model blog.Label
	var count int64
	if err := s.Orm.Debug().Model(&model).Where("id<>? and label_name=?", d.ID, d.LabelName).Count(&count).Error; err != nil {
		s.Log.Errorf("Update Count error:%s", err)
		return err
	}
	if count > 0 {
		return errors.New("标签名称不能重复！")
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

func (s *Label) Insert(d *dto.SaveLabelReq) error {

	var model blog.Label
	var count int64
	if err := s.Orm.Debug().Model(&model).Where("label_name=?", d.LabelName).Count(&count).Error; err != nil {
		s.Log.Errorf("InsertReq Count error:%s", err)
		return err
	}

	if count > 0 {
		return errors.New("标签名称不能重复！")
	}

	data := blog.Label{}
	d.CopyTo(&data)
	err := s.Orm.Debug().Model(&model).Create(&data).Error
	if err != nil {
		s.Log.Errorf("InsertReq error:%s", err)
		return err
	}
	return nil
}

func (s *Label) Delete(d *dto.DelLabelReq) error {
	var err error
	var data blog.Label

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

func (s *Label) Page(d *dto.PageLabelReq, p *actions.DataPermission, page *vo.Page[blog.Label]) error {

	var data blog.Label
	//加载关联表 like
	//.Preload("REProductGroup")
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
