package handler

import (
	"gorm.io/gorm"
	log "wan_go/core/logger"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/utils"
	"wan_go/sdk/pkg"
)

type Login struct {
	UserName string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Captcha  string `form:"captcha" json:"captcha" binding:"required"`
	UUID     string `form:"uuid" json:"uuid" binding:"required"`
	IsAdmin  bool   `form:"isAdmin" json:"isAdmin"`
}

func (u *Login) GetUser(tx *gorm.DB) (user blog.User, role blog.Role, err error) {
	err = tx.Table("user").Where("user_name = ?", u.UserName).First(&user).Error
	if err != nil {
		log.Errorf("get user error, %s", err.Error())
		return
	}
	dec := utils.AesDecryptCrypotJsKey(u.Password)
	_, err = pkg.CompareHashAndPassword(user.Password, dec)
	if err != nil {
		log.Errorf("user login error, %s", err.Error())
		return
	}
	err = tx.Table("role").Where("id = ? ", user.RoleId).First(&role).Error
	if err != nil {
		log.Errorf("get role error, %s", err.Error())
		return
	}
	return
}
