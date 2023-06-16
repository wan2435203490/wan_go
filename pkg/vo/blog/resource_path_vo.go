package blog

import (
	"time"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/sdk/api"
)

// todo time format
type ResourcePathVO struct {
	pagination   *api.Pagination
	ID           int32     `json:"id"`
	Title        string    `json:"title" vd:"@:len($)>0; msg:'资源标题不能为空'"`
	Classify     string    `json:"classify"`
	Cover        string    `json:"cover"`
	Url          string    `json:"url"`
	Type         string    `json:"type" ` //vd:"@:len($)>0; msg:'资源类型不能为空'"
	Remark       string    `json:"remark"`
	Status       bool      `json:"status"`
	Introduction string    `json:"introduction"`
	CreatedAt    time.Time `json:"createAt"`
}

func (to *ResourcePathVO) Copy(from *blog.ResourcePath) {
	to.ID = from.ID
	to.Title = from.Title
	to.Classify = from.Classify
	to.Cover = from.Cover
	to.Url = from.Url
	to.Type = from.Type
	to.Remark = from.Remark
	to.Status = from.Status
	to.Introduction = from.Introduction
	to.CreatedAt = from.CreatedAt
}
