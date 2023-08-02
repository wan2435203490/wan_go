package msg

import (
	"wan_go/pkg/common/constant/landlord_const"
	"wan_go/pkg/common/db/mysql/blog"
)

type playerJoin struct {
	Message
	User *blog.User `json:"user"`
}

func NewPlayerJoin(user *blog.User) *playerJoin {
	var v playerJoin
	v.User = user
	v.Type = v.GetMessageType()
	return &v
}

func (p *playerJoin) GetMessageType() string {
	return landlord_const.PlayerJoin.GetWsMessageType()
}
