package dto

import (
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/models"
	r "wan_go/pkg/common/response"
)

// SaveWeiYanReq update and insert
type SaveWeiYanReq struct {
	ID        int32  `uri:"id"`
	UserId    int32  `json:"user_id"`
	LikeCount int32  `json:"like_count"`
	Content   string `json:"content" vd:"@:len($)>0; msg:'微言不能为空！'"`
	Type      string `json:"type"`
	Source    int32  `json:"source" vd:"$>0; msg:'来源不能为空！'"`
	IsPublic  bool   `json:"is_public"`
	models.ControlBy
	models.ModelTime
}

func (from *SaveWeiYanReq) CopyTo(to *blog.WeiYan) {
	if from.ID != 0 {
		to.ID = from.ID
	}
	from.UserId = to.UserId
	from.LikeCount = to.LikeCount
	from.Content = to.Content
	from.Type = to.Type
	from.Source = to.Source
	from.IsPublic = to.IsPublic
}

func (s *SaveWeiYanReq) GetId() interface{} {
	return s.ID
}

type DelWeiYanReq struct {
	Ids []int `json:"ids"`
}

func (s *DelWeiYanReq) GetId() interface{} {
	return s.Ids
}

type PageWeiYanReq struct {
	*r.Pagination `json:",inline"`
	//ID            int32  `form:"id" search:"type:eq;column:id;table:wei_yan"`
	UserId   int32  `form:"user_id" search:"type:eq;column:user_id;table:wei_yan"`
	Type     string `form:"type" search:"type:eq;column:type;table:wei_yan"`
	IsPublic bool   `form:"is_public" search:"type:eq;column:is_public;table:wei_yan"`
	Source   int32  `form:"source" search:"type:eq;column:source;table:wei_yan"`
}

func (s *PageWeiYanReq) GetNeedSearch() interface{} {
	return *s
}

type PageNewsReq struct {
	*r.Pagination `json:",inline"`
	//ID            int32  `form:"id" search:"type:eq;column:id;table:wei_yan"`
	Type     string `form:"type" search:"type:eq;column:type;table:wei_yan"`
	IsPublic *bool  `form:"isPublic" search:"type:eq;column:is_public;table:wei_yan"`
	Source   int32  `form:"source" search:"type:eq;column:source;table:wei_yan"`
}

func (s *PageNewsReq) GetNeedSearch() interface{} {
	return *s
}
