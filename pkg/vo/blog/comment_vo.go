package blog

import "time"

//todo time format

type CommentVO struct {
	Id     int    `json:"id,omitempty"`
	Source int    `json:"source,omitempty" vd:"$>0; msg:'评论来源标识不能为空'"`
	Type   string `json:"type,omitempty" vd:"@:len($)>0; msg:'评论来源类型不能为空'"`
	//层主的parentCommentId是0，回复的parentCommentId是层主的id
	ParentCommentId int `json:"parentCommentId,omitempty"`
	//层主的parentUserId是null，回复的parentUserId是被回复的userId
	ParentUserId   int    `json:"parentUserId,omitempty"`
	UserId         int    `json:"userId,omitempty"`
	LikeCount      int    `json:"likeCount,omitempty"`
	CommentContent string `json:"commentContent,omitempty" vd:"@:len($)>0; msg:'评论内容不能为空'"`
	CommentInfo    string `json:"commentInfo,omitempty"`
	//子评论必须传评论楼层ID
	FloorCommentId int       `json:"floorCommentId,omitempty"`
	CreateTime     time.Time `json:"createTime"`

	// 需要查询封装 todo
	// childComments Page
	ParentUsername string `json:"parentUsername,omitempty"`
	UserName       string `json:"username,omitempty"`
	Avatar         string `json:"avatar,omitempty"`
}
