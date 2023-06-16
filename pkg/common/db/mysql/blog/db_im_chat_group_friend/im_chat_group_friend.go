package db_im_chat_group_friend

import (
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
)

func Update(data *blog.ImChatUserFriend) error {
	return db.Mysql().Updates(&data).Error
}

func Insert(data *blog.ImChatUserFriend) error {
	return db.Mysql().Create(&data).Error
}

func Delete(data *blog.ImChatUserFriend) error {
	return db.Mysql().Delete(&data).Error
}

func List() (result []*blog.ImChatUserFriend) {

	if err := db.Mysql().
		Find(&result).Error; err != nil {
		return
	}

	return
}
