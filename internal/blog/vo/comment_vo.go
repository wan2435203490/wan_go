package vo

import (
	"time"
	"wan_go/pkg/common/db/mysql/blog"
)

//todo time format

type CommentVO struct {
	ID              int32     `json:"id"`
	Source          int32     `json:"source" ` //vd:"$>0; msg:'评论来源标识不能为空'"
	Type            string    `json:"type" vd:"@:len($)>0; msg:'评论来源类型不能为空'"`
	ParentCommentId int32     `json:"parentCommentId"` //层主的parentCommentId是0，回复的parentCommentId是层主的id
	ParentUserId    int32     `json:"parentUserId"`    //层主的parentUserId是null，回复的parentUserId是被回复的userId
	UserId          int32     `json:"userId"`
	LikeCount       int32     `json:"likeCount"`
	CommentContent  string    `json:"commentContent" vd:"@:len($)>0; msg:'评论内容不能为空'"`
	CommentInfo     string    `json:"commentInfo"`
	FloorCommentId  int32     `json:"floorCommentId"` //子评论必须传评论楼层ID
	CreatedAt       time.Time `json:"createdAt"`
	ChildComments   any       `json:"childComments"` // 需要查询封装 todo
	ParentUserName  string    `json:"parentUserName"`
	UserName        string    `json:"username"`
	Avatar          string    `json:"avatar"`
}

func (to *CommentVO) Copy(from *blog.CommentExt) {
	to.ID = from.ID
	to.Source = from.Source
	to.Type = from.Type
	to.ParentCommentId = from.ParentCommentId
	to.UserId = from.UserId
	to.UserName = from.UserName
	to.FloorCommentId = from.FloorCommentId
	to.ParentUserId = from.ParentUserId
	to.LikeCount = from.LikeCount
	to.CommentContent = from.CommentContent
	to.CommentInfo = from.CommentInfo
	to.CreatedAt = from.CreatedAt
}
