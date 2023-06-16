package db_tree_hole

import (
	"math/rand"
	"time"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
)

func Get(data *blog.TreeHole) error {
	return db.Mysql().Find(&data).Error
}

func Update(data *blog.TreeHole) error {
	return db.Mysql().Updates(&data).Error
}

func Insert(data *blog.TreeHole) error {
	return db.Mysql().Create(&data).Error
}

func Delete(data *blog.TreeHole) error {
	return db.Mysql().Delete(&data).Error
}

func List() (result []*blog.TreeHole) {

	var count64 int64
	if err := db.Mysql().Model(&blog.TreeHole{}).Count(&count64).Error; err != nil {
		return
	}

	count := int(count64)

	var offset int
	if count > blog_const.TREE_HOLE_COUNT {
		r := rand.New(rand.NewSource(time.Millisecond.Milliseconds()))
		offset = r.Intn(count - blog_const.TREE_HOLE_COUNT)
	}

	if err := db.Mysql().
		Offset(offset).
		Limit(blog_const.TREE_HOLE_COUNT).
		Find(&result).Error; err != nil {
		return
	}

	return
}
