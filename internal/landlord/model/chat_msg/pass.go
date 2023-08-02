package msg

import (
	"wan_go/pkg/common/constant/landlord_const"
	"wan_go/pkg/common/db/mysql/blog"
)

type pass struct {
	Message
	User *blog.User `json:"user"`
}

func NewPass(user *blog.User) *pass {
	var v pass
	v.User = user
	v.Type = v.GetMessageType()
	return &v
}

func (p *pass) GetMessageType() string {
	return landlord_const.PassType.GetWsMessageType()
}
