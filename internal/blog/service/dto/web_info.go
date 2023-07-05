package dto

import (
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/models"
)

type SaveWebInfoReq struct {
	ID               int32  `uri:"id"`
	WebName          string `json:"webName,omitempty" vd:"@:len($)>0; msg:'网站名称不能为空！'"`
	WebTitle         string `json:"webTitle,omitempty"`
	Notices          string `json:"notices,omitempty"`
	Footer           string `json:"footer,omitempty"`
	BackgroundImage  string `json:"backgroundImage,omitempty"`
	Avatar           string `json:"avatar,omitempty"`
	RandomAvatar     string `json:"randomAvatar,omitempty"`
	RandomName       string `json:"randomName,omitempty"`
	RandomCover      string `json:"randomCover,omitempty"`
	WaifuJson        string `json:"waifuJson,omitempty"`
	Status           bool   `json:"status,omitempty"`
	models.ControlBy `json:"models.ControlBy"`
}

func (from *SaveWebInfoReq) CopyTo(to *blog.WebInfo) {
	if from.ID != 0 {
		to.ID = from.ID
	}

	to.WebName = from.WebName
	to.WebTitle = from.WebTitle
	to.Notices = from.Notices
	to.Footer = from.Footer
	to.BackgroundImage = from.BackgroundImage
	to.Avatar = from.Avatar
	to.RandomAvatar = from.RandomAvatar
	to.RandomName = from.RandomName
	to.RandomCover = from.RandomCover
	to.WaifuJson = from.WaifuJson
	to.Status = from.Status
}

func (s *SaveWebInfoReq) GetId() interface{} {
	return s.ID
}
