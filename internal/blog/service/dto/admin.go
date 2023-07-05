package dto

import (
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/models"
	r "wan_go/pkg/common/response"
)

// SaveAdminReq update and insert
type SaveAdminReq struct {
	ID               int32  `uri:"id"`
	UserName         string `json:"userName,omitempty"`
	Password         string `json:"password,omitempty"`
	PhoneNumber      string `json:"phoneNumber,omitempty"`
	Email            string `json:"email,omitempty"`
	UserStatus       bool   `json:"userStatus,omitempty"`
	Gender           int8   `json:"gender,omitempty"`
	OpenId           string `json:"openId,omitempty"`
	Avatar           string `json:"avatar,omitempty"`
	Admire           string `json:"admire,omitempty"`
	Introduction     string `json:"introduction,omitempty"`
	UserType         int8   `json:"userType,omitempty"`
	RoleId           int32  `json:"roleId,omitempty"`
	CrypotJsText     string `json:"-"`
	models.ControlBy `json:",inline"`
	models.ModelTime `json:",inline"`
}

func (from *SaveAdminReq) CopyTo(to *blog.User) {
	if from.ID != 0 {
		to.ID = from.ID
	}
	to.ID = from.ID
	to.UserName = from.UserName
	to.Password = from.Password
	to.PhoneNumber = from.PhoneNumber
	to.Email = from.Email
	to.UserStatus = from.UserStatus
	to.Gender = from.Gender
	to.OpenId = from.OpenId
	to.Avatar = from.Avatar
	to.Admire = from.Admire
	to.Introduction = from.Introduction
	to.UserType = from.UserType
	to.RoleId = from.RoleId
	to.CrypotJsText = from.CrypotJsText
}

func (s *SaveAdminReq) GetId() interface{} {
	return s.ID
}

type DelAdminReq struct {
	Ids []int `json:"ids"`
}

func (s *DelAdminReq) GetId() interface{} {
	return s.Ids
}

type PageAdminReq struct {
	*r.Pagination `json:",inline"`
	Account       string `form:"account"`
	UserStatus    *bool  `form:"userStatus" search:"type:eq;column:user_status;table:user"`
	UserType      int    `form:"userType" search:"type:eq;column:user_type;table:user"`
}

func (s *PageAdminReq) GetNeedSearch() interface{} {
	return *s
}

type SaveUserReq struct {
	UserId     int32  `uri:"userId"`
	UserStatus bool   `json:"userStatus"`
	Admire     string `json:"admire"`
	UserType   int    `json:"userType"`
}

func (s *SaveUserReq) GetId() interface{} {
	return s.UserId
}
