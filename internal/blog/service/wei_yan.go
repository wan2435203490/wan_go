package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
	sDto "wan_go/pkg/common/dto"
	"wan_go/sdk/pkg"
	"wan_go/sdk/service"
)

type WeiYan struct {
	service.Service
}

func NewWeiYan(c *gin.Context) *WeiYan {
	ws := WeiYan{}
	ws.Orm = pkg.Orm(c)
	ws.Log = api.GetRequestLogger(c)
	return &ws
}

func (s *WeiYan) InsertReq(d *dto.SaveWeiYanReq) error {

	data := blog.WeiYan{}
	d.CopyTo(&data)

	return s.Insert(&data)
}

func (s *WeiYan) Insert(data *blog.WeiYan) error {

	err := s.Orm.Debug().Model(&data).Create(&data).Error
	if err != nil {
		s.Log.Errorf("InsertReq error:%s", err)
		return err
	}
	return nil
}

func (s *WeiYan) Delete(d *dto.DelWeiYanReq, userId int32) error {
	var err error
	var data blog.WeiYan

	tx := s.Orm.Debug().Model(&data).
		Where("user_id=?", userId).Delete(&data, d.GetId())

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

func (s *WeiYan) Page(d *dto.PageWeiYanReq, p *actions.DataPermission, page *vo.Page[blog.WeiYan], userId int32) error {

	d.Type = blog_const.WEIYAN_TYPE_FRIEND
	//todo ListWeiYan
	if d.UserId == 0 {
		d.UserId = userId
	}

	var data blog.WeiYan
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

// todo
func ListWeiYan(vo *vo.BaseRequestVO[*blog.WeiYan], userId int32, adminId int32) {

	tx := db.Page(&vo.Pagination).Where("type=?", blog_const.WEIYAN_TYPE_FRIEND)

	if vo.UserId == 0 {
		if userId > 0 {
			tx.Where("user_id=?", userId)
		} else {
			tx.Where("is_public=? and user_id=?", blog_const.PUBLIC.Code, adminId)
		}
	} else {
		if vo.UserId != int32(userId) {
			tx.Where("is_public=?", blog_const.PUBLIC.Code)
		}
		tx.Where("user_id=?", vo.UserId)
	}

	tx.Order("created_at DESC").Find(&vo.Records)
}

func (s *WeiYan) PageNews(d *dto.PageNewsReq, p *actions.DataPermission, page *vo.Page[blog.WeiYan]) error {

	d.Type = blog_const.WEIYAN_TYPE_NEWS
	isPublic := blog_const.PUBLIC.Code == 1
	d.IsPublic = &isPublic

	var data blog.WeiYan
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
