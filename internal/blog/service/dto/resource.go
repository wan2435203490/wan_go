package dto

import (
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/models"
	r "wan_go/pkg/common/response"
)

// SaveResourceReq update and insert
type SaveResourceReq struct {
	ID               int32  `uri:"id"`
	UserId           int32  `json:"userId,omitempty"`
	Type             string `json:"type,omitempty" vd:"@:len($)>0; msg:'资源类型不能为空！'"`
	Path             string `json:"path,omitempty" vd:"@:len($)>0; msg:'资源路径不能为空！'"`
	Size             int32  `json:"size,omitempty"`
	MimeType         string `json:"mimeType,omitempty"`
	Status           bool   `json:"status,omitempty"`
	models.ControlBy `json:"models.ControlBy"`
}

func (from *SaveResourceReq) CopyTo(to *blog.Resource) {
	if from.ID != 0 {
		to.ID = from.ID
	}
	to.UserId = from.UserId
	to.Type = from.Type
	to.Path = from.Path
	to.Size = from.Size
	to.MimeType = from.MimeType
	to.Status = from.Status
}

func (s *SaveResourceReq) GetId() interface{} {
	return s.ID
}

type DelResourceReq struct {
	Path string `json:"path"`
}

func (s *DelResourceReq) GetPath() interface{} {
	return s.Path
}

type PageResourceReq struct {
	*r.Pagination `json:",inline"`
	ID            int32  `uri:"id" search:"type:eq;column:id;table:resource" comment:"id"`
	ResourceType  string `form:"resourceType" search:"type:eq;column:type;table:resource"`
}

func (s *PageResourceReq) GetNeedSearch() interface{} {
	return *s
}

type ChangeResourceReq struct {
	ID     int32 `uri:"ID"`
	Status *bool `json:"status" vd:"$==nil; msg='status is nil'"`
}
