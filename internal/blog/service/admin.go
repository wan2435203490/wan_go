package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/db/mysql/blog"
	sDto "wan_go/pkg/common/dto"
	"wan_go/sdk/service"
)

type Admin struct {
	service.Service
}

func (s *Admin) Update(d *dto.SaveAdminReq) error {

	var model blog.User
	var count int64
	if err := s.Orm.Debug().Model(&model).Where("id<>? and sort_name=?", d.ID, d.UserName).Count(&count).Error; err != nil {
		s.Log.Errorf("Update Count error:%s", err)
		return err
	}
	if count > 0 {
		return errors.New("用户名称不能重复！")
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

func (s *Admin) Insert(d *dto.SaveAdminReq) error {

	var model blog.User
	var count int64
	if err := s.Orm.Debug().Model(&model).Where("sort_name=?", d.UserName).Count(&count).Error; err != nil {
		s.Log.Errorf("InsertReq Count error:%s", err)
		return err
	}

	if count > 0 {
		return errors.New("用户名称不能重复！")
	}

	data := blog.User{}
	d.CopyTo(&data)
	err := s.Orm.Debug().Model(&model).Create(&data).Error
	if err != nil {
		s.Log.Errorf("InsertReq error:%s", err)
		return err
	}
	return nil
}

func (s *Admin) Delete(d *dto.DelAdminReq) error {
	var err error
	var data blog.User

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

func (s *Admin) Page(d *dto.PageAdminReq, p *actions.DataPermission, page *vo.Page[blog.User]) error {

	var data blog.User
	err := s.Orm.Debug().Scopes(
		sDto.MakeCondition(d.GetNeedSearch()),
		func(db *gorm.DB) *gorm.DB {
			account := fmt.Sprintf("%%%s%%", d.Account)
			return db.Where("user_name like ? or email like ? or phone_number like ?", account, account, account)
		},
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

func (s *Admin) UpdateUserStatus(userId int32, userStatus bool) error {
	tx := s.Orm.Debug().Model(&blog.User{ID: userId}).Where("user_status=?", !userStatus).
		Update("user_status", userStatus)
	if err := tx.Error; err != nil {
		s.Log.Errorf("Update error:%s", err)
		return err
	}
	if tx.RowsAffected == 0 {
		return errors.New("数据不存在或无权更新该数据")
	}
	return nil
}

func (s *Admin) UpdateUserAdmire(userId int32, admire string) error {
	tx := s.Orm.Debug().Model(&blog.User{ID: userId}).
		Update("admire", admire)
	if err := tx.Error; err != nil {
		s.Log.Errorf("Update error:%s", err)
		return err
	}
	if tx.RowsAffected == 0 {
		return errors.New("数据不存在或无权更新该数据")
	}
	return nil
}

func (s *Admin) UpdateUserType(userId int32, userType int) error {
	tx := s.Orm.Debug().Model(&blog.User{ID: userId}).
		Update("user_type", userType)
	if err := tx.Error; err != nil {
		s.Log.Errorf("Update error:%s", err)
		return err
	}
	if tx.RowsAffected == 0 {
		return errors.New("数据不存在或无权更新该数据")
	}
	return nil
}
