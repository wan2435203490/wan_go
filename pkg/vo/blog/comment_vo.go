package blog

import (
	"time"
	"wan_go/pkg/common/db/mysql/blog"
)

//todo time format

type CommentVO struct {
	ID     int32  `json:"id,omitempty"`
	Source int32  `json:"source,omitempty" vd:"$>0; msg:'评论来源标识不能为空'"`
	Type   string `json:"type,omitempty" vd:"@:len($)>0; msg:'评论来源类型不能为空'"`
	//层主的parentCommentId是0，回复的parentCommentId是层主的id
	ParentCommentId int32 `json:"parentCommentId,omitempty"`
	//层主的parentUserId是null，回复的parentUserId是被回复的userId
	ParentUserId   int32  `json:"parentUserId,omitempty"`
	UserId         int32  `json:"userId,omitempty"`
	LikeCount      int32  `json:"likeCount,omitempty"`
	CommentContent string `json:"commentContent,omitempty" vd:"@:len($)>0; msg:'评论内容不能为空'"`
	CommentInfo    string `json:"commentInfo,omitempty"`
	//子评论必须传评论楼层ID
	FloorCommentId int32     `json:"floorCommentId,omitempty"`
	CreatedAt      time.Time `json:"createTime"`

	// 需要查询封装 todo
	//ChildComments  api.Pagination
	ParentUsername string `json:"parentUsername,omitempty"`
	UserName       string `json:"username,omitempty"`
	Avatar         string `json:"avatar,omitempty"`
}

func (to *CommentVO) Copy(from *blog.Comment) {
	to.ID = from.ID
	to.Source = from.Source
	to.Type = from.Type
	to.ParentCommentId = from.ParentCommentId
	to.UserId = from.UserId
	to.FloorCommentId = from.FloorCommentId
	to.ParentUserId = from.ParentUserId
	to.LikeCount = from.LikeCount
	to.CommentContent = from.CommentContent
	to.CommentInfo = from.CommentInfo
	to.CreatedAt = from.CreatedAt
}
