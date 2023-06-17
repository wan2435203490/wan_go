package db_family

import (
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
	blogVO "wan_go/pkg/vo/blog"
)

func GetByUserId(userId int) *blog.Family {
	family := blog.Family{}
	if db.Mysql().Where("user_id=?", userId).First(&family).Error != nil {
		return nil
	}
	return &family
}

func Update(family *blog.Family) {
	db.Mysql().UpdateColumns(&family)
}

func Insert(family *blog.Family) {
	db.Mysql().Create(&family)
}

func DeleteById(id int) {
	db.Mysql().Delete(&blog.Family{ID: int32(id)})
}

func ListFamily(vo *blogVO.BaseRequestVO[*blog.Family]) {
	db.Page(&vo.Pagination).Where("status=?", vo.Status).
		Order("CreatedAt DESC").Find(&vo.Records)
}

func ChangeLoveStatus(id int, status bool) {
	db.Mysql().Model(&blog.Family{ID: int32(id)}).Update("status", status)
}
