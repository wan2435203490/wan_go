package dto

import (
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/models"
	r "wan_go/pkg/common/response"
)

// SaveLabelReq update and insert
type SaveLabelReq struct {
	ID               int32  `uri:"id"`
	SortId           int32  `json:"sortId" vd:"$>0; msg:'分类Id不能为空！'"`
	LabelName        string `json:"labelName" vd:"@:len($)>0; msg:'标签名称不能为空！'"`
	LabelDescription string `json:"labelDescription" vd:"@:len($)>0; msg:'标签描述不能为空！'"`
	models.ControlBy
}

func (from *SaveLabelReq) CopyTo(to *blog.Label) {
	if from.ID != 0 {
		to.ID = from.ID
	}
	to.SortId = from.SortId
	to.LabelName = from.LabelName
	to.LabelDescription = from.LabelDescription
}

func (s *SaveLabelReq) GetId() interface{} {
	return s.ID
}

type DelLabelReq struct {
	Ids []int `json:"ids"`
}

func (s *DelLabelReq) GetId() interface{} {
	return s.Ids
}

type PageLabelReq struct {
	*r.Pagination    `json:",inline"`
	ID               int32  `uri:"id" search:"type:eq;column:id;table:label" comment:"id"`
	SortId           int32  `form:"sortId" search:"type:eq;column:sort_id;table:label" comment:"分类id"`
	LabelName        string `form:"labelName" search:"type:like;column:label_name;table:label" comment:"标签名称"`
	LabelDescription string `form:"labelDescription" search:"type:like;column:label_description;table:label" comment:"标签描述"`
}

func (s *PageLabelReq) GetNeedSearch() interface{} {
	return *s
}
