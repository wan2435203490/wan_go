package blog

import (
	"time"
	"wan_go/pkg/common/db/mysql/blog"
)

type ArticleVO struct {
	ID     int `json:"id,omitempty"`
	UserId int `json:"userId,omitempty"`
	// 查询为空时，随机选择
	ArticleCover    string    `json:"articleCover,omitempty"`
	ArticleTitle    string    `json:"articleTitle,omitempty" vd:"@:len($)>0; msg:'用户名不能为空'"`
	ArticleContent  string    `json:"articleContent,omitempty" vd:"@:len($)>0; msg:'用户名不能为空'"`
	ViewCount       int       `json:"viewCount,omitempty"`
	LikeCount       int       `json:"likeCount,omitempty"`
	CommentStatus   bool      `json:"commentStatus,omitempty"`
	RecommendStatus bool      `json:"recommendStatus,omitempty"`
	Password        string    `json:"password,omitempty"`
	ViewStatus      bool      `json:"viewStatus,omitempty"`
	CreateTime      time.Time `json:"createTime"`
	UpdateTime      time.Time `json:"updateTime"`
	UpdateBy        string    `json:"updateBy,omitempty"`
	SortId          int       `json:"sortId,omitempty" vd:"$>0; msg:'文章分类不能为空'"`
	LabelId         int       `json:"labelId,omitempty" vd:"$>0; msg:'文章标签不能为空'"`
	// 需要查询封装
	CommentCount int         `json:"commentCount,omitempty"`
	Username     string      `json:"username,omitempty"`
	Sort         *blog.Sort  `json:"sort,omitempty"`
	Label        *blog.Label `json:"label,omitempty"`
}
