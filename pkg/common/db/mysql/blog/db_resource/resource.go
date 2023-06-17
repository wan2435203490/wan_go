package db_resource

import (
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
)

func Insert(data *blog.Resource) error {
	return db.Mysql().Create(&data).Error
}

func DeleteByPath(path string) error {
	return db.Mysql().
		Where("path=?", path).
		Delete(&blog.Resource{}).Error
}

func DeleteByUserId(id int) error {
	userId := cache.GetUserId()
	return db.Mysql().
		Where("id=? and user_id=?", id, userId).
		Delete(&blog.Resource{}).Error
}

func GetResourceInfo() *[]*blog.Resource {
	var resources []*blog.Resource
	db.Mysql().Select("id, path").
		Where("path like ? and size is null", "%"+blog_const.DOWNLOAD_URL+"%").
		Find(&resources)

	return &resources
}

func BatchUpdate(resources *[]*blog.Resource) error {
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

func GetImageList() *[]*blog.Resource {
	var resources []*blog.Resource
	db.Mysql().Where("type=? and status=? and user_id=?",
		blog_const.PATH_TYPE_INTERNET_MEME, blog_const.STATUS_ENABLE.Code, cache.GetAdminUserId()).
		Order("CreatedAt DESC").Find(&resources)
	return &resources
}

func ListResource(vo *blogVO.BaseRequestVO[*blog.Resource]) {

	tx := db.Page(&vo.Pagination)

	if utils.IsNotEmpty(vo.ResourceType) {
		tx.Where("type=?", vo.ResourceType)
	}

	tx.Order("CreatedAt DESC").Find(&vo.Records)
}

func ChangeResourceStatus(id int, status bool) {
	db.Mysql().Model(&blog.Resource{ID: int32(id)}).Update("status", status)
}
