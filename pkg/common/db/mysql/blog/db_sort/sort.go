package db_sort

import (
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
)

func Update(data *blog.Sort) error {
	return db.Mysql().Updates(&data).Error
}

func Insert(data *blog.Sort) error {
	return db.Mysql().Create(&data).Error
}

func Delete(data *blog.Sort) error {
	return db.Mysql().Delete(&data).Error
}

func List() (result []*blog.Sort) {

	if err := db.Mysql().
		Find(&result).Error; err != nil {
		return
	}

	return
}
