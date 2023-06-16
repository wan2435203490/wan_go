package db_im_chat_group_user

import (
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
)

func Update(data *blog.ImChatGroupUser) error {
	return db.Mysql().Updates(&data).Error
}

func Insert(data *blog.ImChatGroupUser) error {
	return db.Mysql().Create(&data).Error
}

func Delete(data *blog.ImChatGroupUser) error {
	return db.Mysql().Delete(&data).Error
}

func List() (result []*blog.ImChatGroupUser) {

	if err := db.Mysql().
		Find(&result).Error; err != nil {
		return
	}

	return
}
