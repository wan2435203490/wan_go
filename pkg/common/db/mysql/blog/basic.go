package blog

import (
	"gorm.io/gorm"
	"wan_go/pkg/common/db"
)

func Where(query interface{}, args ...interface{}) *gorm.DB {
	return db.DB.MysqlDB.Where(query, args)
}

func Select(query interface{}, args ...interface{}) *gorm.DB {
	return db.DB.MysqlDB.Select(query, args)
}
