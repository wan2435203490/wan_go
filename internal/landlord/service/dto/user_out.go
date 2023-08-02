package dto

import (
	"wan_go/pkg/common/db/mysql/blog"
)

type UserOut struct {
	ID       int32   `json:"id"`
	UserName string  `json:"username"`
	Gender   int8    `json:"gender"`
	Avatar   string  `json:"avatar"`
	Money    float64 `json:"money"`
}

func ToUserOut(u *blog.User) *UserOut {
	if u == nil {
		return nil
	}

	return &UserOut{
		ID:       u.ID,
		UserName: u.UserName,
		Avatar:   u.Avatar,
		Gender:   u.Gender,
		Money:    u.Money,
	}
}

func ToUserOutList(users []*blog.User) []*UserOut {
	if users == nil {
		return nil
	}

	var ret []*UserOut

	for _, u := range users {
		user := ToUserOut(u)
		ret = append(ret, user)
	}

	return ret
}
