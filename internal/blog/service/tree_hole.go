package service

import (
	"errors"
	"math/rand"
	"time"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"
	sDto "wan_go/pkg/common/dto"
	"wan_go/sdk/service"
)

type TreeHole struct {
	service.Service
}

func (s *TreeHole) Insert(d *dto.SaveTreeHoleReq) error {

	data := blog.TreeHole{}
	d.CopyTo(&data)
	err := s.Orm.Debug().Model(&data).Create(&data).Error
	if err != nil {
		s.Log.Errorf("InsertReq error:%s", err)
		return err
	}
	return nil
}

func (s *TreeHole) Delete(d *dto.DelTreeHoleReq) error {
	var err error
	var data blog.TreeHole

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

func (s *TreeHole) PageUser(d *dto.PageTreeHoleReq, p *actions.DataPermission, page *vo.Page[blog.TreeHole]) error {

	var count64 int64
	if err := s.Orm.Model(&blog.TreeHole{}).Count(&count64).Error; err != nil {
		return err
	}

	count := int(count64)
	var offset int
	if count > blog_const.TREE_HOLE_COUNT {
		r := rand.New(rand.NewSource(time.Millisecond.Milliseconds()))
		offset = r.Intn(count - blog_const.TREE_HOLE_COUNT)
	}
	d.Current = offset + 1
	d.Size = blog_const.TREE_HOLE_COUNT

	return s.Page(d, p, page)
}

func (s *TreeHole) Page(d *dto.PageTreeHoleReq, p *actions.DataPermission, page *vo.Page[blog.TreeHole]) error {

	var data blog.TreeHole
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
