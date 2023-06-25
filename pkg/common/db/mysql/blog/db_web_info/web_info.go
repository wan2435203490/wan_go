package db_web_info

import (
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
)

func Update(webInfo *blog.WebInfo) error {
	return db.DB.MysqlDB.Updates(&webInfo).Error
}

func List() ([]*blog.WebInfo, error) {
	var webInfoes []*blog.WebInfo

	if err := db.DB.MysqlDB.Find(&webInfoes).Error; err != nil {
		return nil, err
	}
	return webInfoes, nil
}
