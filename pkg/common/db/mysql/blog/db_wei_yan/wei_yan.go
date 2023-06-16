package db_wei_yan

import (
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
)

func Update(data *blog.WeiYan) error {
	return db.Mysql().Updates(&data).Error
}

func Insert(data *blog.WeiYan) error {
	return db.Mysql().Create(&data).Error
}

func Delete(data *blog.WeiYan) error {
	return db.Mysql().Delete(&data).Error
}

func List() (result []*blog.WeiYan) {

	if err := db.Mysql().
		Find(&result).Error; err != nil {
		return
	}

	return
}
