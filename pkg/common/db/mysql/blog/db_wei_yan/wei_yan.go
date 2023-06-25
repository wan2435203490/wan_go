package db_wei_yan

import (
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
	blogVO "wan_go/pkg/vo/blog"
)

func Insert(data *blog.WeiYan) error {
	return db.Mysql().Create(&data).Error
}

func DeleteByUserId(id int) error {
	userId := cache.GetUserId()
	return db.Mysql().
		Where("id=? and user_id=?", id, userId).
		Delete(&blog.WeiYan{}).Error
}

func ListNews(vo *blogVO.BaseRequestVO[*blog.WeiYan]) {

	db.Page(&vo.Pagination).Where("type=? and source=? and is_public=?",
		blog_const.WEIYAN_TYPE_NEWS, vo.Source, blog_const.PUBLIC.Code).
		Order("created_at DESC").Find(&vo.Records)
}

func ListWeiYan(vo *blogVO.BaseRequestVO[*blog.WeiYan]) {

	tx := db.Page(&vo.Pagination).Where("type=?", blog_const.WEIYAN_TYPE_FRIEND)

	userId := cache.GetUserId()

	if vo.UserId == 0 {
		if userId > 0 {
			tx.Where("user_id=?", userId)
		} else {
			userId = cache.GetAdminUserId()
			tx.Where("is_public=? and user_id=?", blog_const.PUBLIC.Code, userId)
		}
	} else {
		if vo.UserId != int32(userId) {
			tx.Where("is_public=?", blog_const.PUBLIC.Code)
		}
		tx.Where("user_id=?", vo.UserId)
	}

	tx.Order("created_at DESC").Find(&vo.Records)
}
