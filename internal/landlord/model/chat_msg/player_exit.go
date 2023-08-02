package msg

import (
	"wan_go/pkg/common/constant/landlord_const"
	"wan_go/pkg/common/db/mysql/blog"
)

type playerExit struct {
	Message
	User *blog.User `json:"user"`
}

func NewPlayerExit(user *blog.User) *playerExit {
	var v playerExit
	v.User = user
	v.Type = v.GetMessageType()
	return &v
}

func (p *playerExit) GetMessageType() string {
	return landlord_const.PlayerExit.GetWsMessageType()
}
