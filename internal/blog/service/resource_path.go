package service

import (
	"errors"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
	sDto "wan_go/pkg/common/dto"
	"wan_go/sdk/service"
)

type ResourcePath struct {
	service.Service
}

func (s *ResourcePath) Insert(d *dto.SaveResourcePathReq) error {

	var data blog.ResourcePath

	d.CopyTo(&data)
	err := s.Orm.Debug().Model(&data).Create(&data).Error
	if err != nil {
		s.Log.Errorf("InsertReq error:%s", err)
		return err
	}
	return nil
}

func (s *ResourcePath) Delete(d *dto.DelResourcePathReq) error {
	var err error
	var data blog.ResourcePath

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

func (s *ResourcePath) Update(d *dto.SaveResourcePathReq) error {

	var model blog.ResourcePath
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

func (s *ResourcePath) ListByResourceTypeAndClassify(d *dto.PageResourcePathReq, p *actions.DataPermission, page *vo.Page[vo.ResourcePathVO]) error {

	var model blog.ResourcePath
	var records []blog.ResourcePath

	if err := s.Orm.Debug().Model(&model).Scopes(
		sDto.MakeCondition(d.GetNeedSearch()),
		sDto.Paginate(d.Pagination),
		actions.Permission(model.TableName(), p),
	).Find(&records).
		Limit(-1).
		Offset(-1).
		Count(&page.Total).Error; err != nil {
		s.Log.Errorf("ListByResourceTypeAndClassify error: %s", err)
		return err
	}

	for _, path := range records {
		pp := vo.ResourcePathVO{}
		pp.Copy(&path)
		page.Records = append(page.Records, pp)
	}

	return nil
}

func (s *ResourcePath) ListFunny() (ret []map[string]any) {
	if err := db.Mysql().Model(&blog.ResourcePath{}).
		Select("classify, count(*) as count").
		Where("status = ? and type = ?", true, blog_const.RESOURCE_PATH_TYPE_FUNNY).
		Group("classify").
		Find(&ret).Error; err != nil {
		return
	}

	return
}

// todo test
func (s *ResourcePath) ListCollect(classifyMap *map[string][]vo.ResourcePathVO) error {

	maps := make(map[string][]vo.ResourcePathVO, 0)
	var paths []*blog.ResourcePath
	if err := s.Orm.Debug().
		Where("status = ? and type = ?", true, blog_const.RESOURCE_PATH_TYPE_FAVORITES).
		Order("classify, title").
		Find(&paths).Error; err != nil {
		s.Log.Errorf("ListCollect error: %s", err)
		return err
	}

	for _, path := range paths {
		var pathVO vo.ResourcePathVO
		pathVO.Copy(path)
		maps[path.Classify] = append(maps[path.Classify], pathVO)
	}

	*classifyMap = maps

	return nil
}

func (s *ResourcePath) ListAdminLovePhoto(adminId int, ret *[]map[string]any) error {
	//todo why it is []map?
	if err := s.Orm.Debug().Model(&blog.ResourcePath{}).
		Select("classify, count(1) as count").
		Where("status = ? and type = ? and remark = ?",
			true, blog_const.RESOURCE_PATH_TYPE_LOVE_PHOTO, adminId).
		Group("classify").
		Find(&ret).Error; err != nil {
		s.Log.Errorf("ListCollect error: %s", err)
		return err
	}

	return nil
}
