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

//
//func GetOrInsertUser(db_user *blog.User) error {
//
//	err := db.MySQL.First(&db_user, "username=?", db_user.UserName).Error
//
//	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
//		//第一次进来 用户不存在 自动插入
//		db_user.Timestamp = time.Now()
//		db_user.Id = strings.ReplaceAll(uuid.NewString(), "-", "")
//		err = db.MySQL.Create(&db_user).Error
//	}
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
