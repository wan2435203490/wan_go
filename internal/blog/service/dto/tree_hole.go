package dto

import (
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/models"
	r "wan_go/pkg/common/response"
)

// SaveTreeHoleReq update and insert
type SaveTreeHoleReq struct {
	ID      int32  `uri:"id"`
	Avatar  string `json:"avatar"`
	Message string `json:"message" vd:"@:len($)>0; msg:'留言不能为空！'"`
	models.ControlBy
}

func (from *SaveTreeHoleReq) CopyTo(to *blog.TreeHole) {
	if from.ID != 0 {
		to.ID = from.ID
	}
	to.Avatar = from.Avatar
	to.Message = from.Message
}

func (s *SaveTreeHoleReq) GetId() interface{} {
	return s.ID
}

type DelTreeHoleReq struct {
	Ids []int `json:"ids"`
}

func (s *DelTreeHoleReq) GetId() interface{} {
	return s.Ids
}

type PageTreeHoleReq struct {
	*r.Pagination `json:",inline"`
	ID            int32  `uri:"id" search:"type:eq;column:id;table:tree_hole" comment:"id"`
	Avatar        string `form:"avatar" search:"type:like;column:avatar;table:tree_hole" comment:"留言头像"`
	Message       string `form:"message" search:"type:like;column:message;table:tree_hole" comment:"留言内容"`
}

func (s *PageTreeHoleReq) GetNeedSearch() interface{} {
	return *s
}
