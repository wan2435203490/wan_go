package db

import (
	"fmt"
	"gorm.io/gorm/schema"
	"time"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/db/mysql/blog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Writer struct{}

func (w Writer) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

const dsnFormat = "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local"

func InitMysql(DB *DataBases) {
	dsn := fmt.Sprintf(dsnFormat,
		config.Config.Mysql.UserName,
		config.Config.Mysql.Password,
		config.Config.Mysql.Address[0],
		"mysql")

	var db *gorm.DB
	var err1 error
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		time.Sleep(time.Duration(30) * time.Second)
		db, err1 = gorm.Open(mysql.Open(dsn), nil)
		if err1 != nil {
			panic(err1.Error() + " open failed " + dsn)
		}
	}

	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8 COLLATE utf8_general_ci;", config.Config.Mysql.DatabaseName)
	err = db.Exec(sql).Error
	if err != nil {
		panic(err.Error() + " Exec failed " + sql)
	}
	dsn = fmt.Sprintf(dsnFormat,
		config.Config.Mysql.UserName,
		config.Config.Mysql.Password,
		config.Config.Mysql.Address[0],
		config.Config.Mysql.DatabaseName)
	newLogger := logger.New(
		Writer{},
		logger.Config{
			SlowThreshold:             time.Duration(config.Config.Mysql.SlowThreshold) * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.LogLevel(config.Config.Mysql.LogLevel),                       // Log level
			IgnoreRecordNotFoundError: true,                                                                // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                                                // Disable color
		},
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		panic(err.Error() + " Open failed " + dsn)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error() + " db.DB() failed ")
	}

	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.Config.Mysql.MaxLifeTime))
	sqlDB.SetMaxOpenConns(config.Config.Mysql.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Config.Mysql.MaxIdleConns)

	db.AutoMigrate(
		&blog.User{},
		&blog.Role{},
		&blog.Article{},
		&blog.Comment{},
		&blog.Sort{},
		&blog.Label{},
		&blog.TreeHole{},
		&blog.WeiYan{},
		&blog.WebInfo{},
		&blog.ResourcePath{},
		&blog.Resource{},
		&blog.Family{},
		&blog.ImChatUserFriend{},
		&blog.ImChatGroup{},
		&blog.ImChatGroupUser{},
		&blog.ImChatUserMessage{},
		&blog.ImChatUserGroupMessage{},
		&blog.Achievement{},
	)

	db.Set("gorm:table_options", "CHARSET=utf8")
	db.Set("gorm:table_options", "collation=utf8_unicode_ci")

	//if !db.Migrator().HasTable(&Friend{}) {
	//	db.Migrator().CreateTable(&Friend{})
	//}

	DB.MysqlDB = db
}
