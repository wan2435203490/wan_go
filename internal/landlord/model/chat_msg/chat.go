package msg

import (
	"time"
	"wan_go/internal/blog/vo"
	"wan_go/internal/landlord/service/dto"
	"wan_go/pkg/common/constant/landlord_const"
	"wan_go/pkg/common/db/mysql/blog"
)

type Chat struct {
	Message
	Sender     *vo.UserVO `json:"sender"`
	Content    string     `json:"content"`
	TypeId     int        `json:"typeId"`
	Dimension  string     `json:"dimension"`
	CreateTime time.Time  `json:"createTime"`
}

func NewChat(c *dto.Chat, user *blog.User) *Chat {
	v := &Chat{
		Content:    c.Content,
		TypeId:     c.Type,
		Sender:     &vo.UserVO{UserName: user.UserName, Avatar: user.Avatar},
		Dimension:  c.Dimension,
		CreateTime: time.Now(),
	}
	v.Type = v.GetMessageType()
	return v
}

func (c *Chat) GetMessageType() string {
	return landlord_const.Chat.GetWsMessageType()
}
