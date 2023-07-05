package dto

import (
	"time"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/models"
	r "wan_go/pkg/common/response"
)

// SaveCommentReq update and insert
type SaveCommentReq struct {
	ID               int32     `uri:"ID,omitempty"`
	Source           int32     `json:"source,omitempty"`
	Type             string    `json:"type,omitempty"`
	ParentCommentId  int32     `json:"parentCommentId,omitempty"`
	ParentUserId     int32     `json:"parentUserId,omitempty"`
	UserId           int32     `json:"userId,omitempty"`
	LikeCount        int32     `json:"likeCount,omitempty"`
	CommentContent   string    `json:"commentContent,omitempty"`
	CommentInfo      string    `json:"commentInfo,omitempty"`
	FloorCommentId   int32     `json:"floorCommentId,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	ChildComments    any       `json:"childComments,omitempty"`
	ParentUsername   string    `json:"parentUsername,omitempty"`
	UserName         string    `json:"userName,omitempty"`
	Avatar           string    `json:"avatar,omitempty"`
	models.ControlBy `json:",inline" json:"models.ControlBy"`
	models.ModelTime `json:",inline" json:"models.ModelTime"`
}

func (from *SaveCommentReq) CopyTo(to *blog.Comment) {
	if from.ID != 0 {
		to.ID = from.ID
	}
	to.Source = from.Source
	to.Type = from.Type
	to.ParentCommentId = from.ParentCommentId
	to.ParentUserId = from.ParentUserId
	to.UserId = from.UserId
	to.LikeCount = from.LikeCount
	to.CommentContent = from.CommentContent
	to.CommentInfo = from.CommentInfo
	to.FloorCommentId = from.FloorCommentId
	to.CreatedAt = from.CreatedAt
	//to.ChildComments     = from.ChildComments
	//to.ParentUsername    = from.ParentUsername
	//to.UserName          = from.UserName
	//to.Avatar            = from.Avatar
}

func (s *SaveCommentReq) GetId() interface{} {
	return s.ID
}

type DelCommentReq struct {
	Ids []int `json:"ids"`
}

func (s *DelCommentReq) GetId() interface{} {
	return s.Ids
}

type PageAdminCommentReq struct {
	*r.Pagination `json:",inline"`
	Source        int32  `form:"source" search:"type:eq;column:source;table:comment"`
	Sources       *[]int `search:"type:in;column:source;table:comment"`
	CommentType   string `form:"commentType" search:"type:eq;column:type;table:comment"`
}

func (s *PageAdminCommentReq) GetNeedSearch() interface{} {
	return *s
}

type PageCommentReq struct {
	*r.Pagination  `json:",inline"`
	Source         int32  `form:"source" vd:"$>-1;msg:'source不能为空！'" search:"type:eq;column:source;table:comment"`
	FloorCommentId int    `form:"floorCommentId" search:"type:eq;column:floor_comment_id;table:comment"`
	CommentType    string `form:"commentType" vd:"@:len($)>0;msg:'commentType不能为空！'" search:"type:eq;column:type;table:comment"`
}

func (s *PageCommentReq) GetNeedSearch() interface{} {
	return *s
}

type CountCommentReq struct {
	Source int32  `form:"source,omitempty"`
	Type   string `form:"type,omitempty"`
}
