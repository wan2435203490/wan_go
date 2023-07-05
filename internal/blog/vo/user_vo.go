package blog

import (
	"time"
	"wan_go/pkg/common/db/mysql/blog"
	r "wan_go/pkg/common/response"
)

// todo 时间格式化
type UserVO struct {
	r.CodeMsg
	ID           int32     `json:"id"`
	UserName     string    `json:"username" ` //vd:"@:len($)>0; msg:'用户名不能为空'"
	Password     string    `json:"password"`
	PhoneNumber  string    `json:"phoneNumber"`
	Email        string    `json:"email"`
	Gender       int8      `json:"gender"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	CreatedAt    time.Time `json:"createTime"`
	UpdatedAt    time.Time `json:"updateTime"`
	UpdateBy     string    `json:"updateBy"`
	IsBoss       bool      `json:"isBoss"`
	AccessToken  string    `json:"accessToken"`
	Code         string    `json:"code"`
}

func (to *UserVO) Copy(from *blog.User) {
	to.ID = from.ID
	to.UserName = from.UserName
	to.Password = from.Password
	to.PhoneNumber = from.PhoneNumber
	to.Email = from.Email
	to.Gender = from.Gender
	to.Avatar = from.Avatar
	to.Introduction = from.Introduction
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
	to.UpdateBy = from.UpdateBy
}
