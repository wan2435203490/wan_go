package dto

import (
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/models"
	r "wan_go/pkg/common/response"
)

// SaveResourcePathReq update and insert
type SaveResourcePathReq struct {
	ID               int32  `uri:"id" `
	Title            string `json:"title,omitempty"`
	Classify         string `json:"classify,omitempty"`
	Cover            string `json:"cover,omitempty"`
	Url              string `json:"url,omitempty"`
	Introduction     string `json:"introduction,omitempty"`
	Type             string `json:"type,omitempty"` //vd:"@:len($)>0; msg:'资源类型不能为空！'"
	Status           bool   `json:"status,omitempty"`
	Remark           string `json:"remark,omitempty"`
	models.ControlBy `json:"models.ControlBy"`
}

func (from *SaveResourcePathReq) CopyTo(to *blog.ResourcePath) {
	if from.ID != 0 {
		to.ID = from.ID
	}
	to.Title = from.Title
	to.Classify = from.Classify
	to.Cover = from.Cover
	to.Url = from.Url
	to.Introduction = from.Introduction
	to.Type = from.Type
	to.Status = from.Status
	to.Remark = from.Remark
}

func (s *SaveResourcePathReq) GetId() interface{} {
	return s.ID
}

type DelResourcePathReq struct {
	Ids []int `json:"ids"`
}

func (s *DelResourcePathReq) GetId() interface{} {
	return s.Ids
}

type PageResourcePathReq struct {
	*r.Pagination `json:",inline"`
	ID            int32  `uri:"id" search:"type:eq;column:id;table:resource_path" comment:"id"`
	Status        *bool  `form:"status" search:"type:eq;column:status;table:resource_path"`
	ResourceType  string `form:"resourceType" search:"type:eq;column:type;table:resource_path"`
	Classify      string `form:"classify" search:"type:eq;column:classify;table:resource_path"`
}

func (s *PageResourcePathReq) GetNeedSearch() interface{} {
	return *s
}
