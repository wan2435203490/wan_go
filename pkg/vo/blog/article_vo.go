package blog

import (
	"time"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/sdk/api"
)

type ArticleVO struct {
	api.CodeMsg
	ID     int32 `json:"id,omitempty"`
	UserId int32 `json:"userId,omitempty"`
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
	CreatedAt       time.Time `json:"createTime"`
	UpdatedAt       time.Time `json:"updateTime"`
	UpdateBy        string    `json:"updateBy,omitempty"`
	SortId          int32     `json:"sortId,omitempty" vd:"$>0; msg:'文章分类不能为空'"`
	LabelId         int32     `json:"labelId,omitempty" vd:"$>0; msg:'文章标签不能为空'"`
	// 需要查询封装
	CommentCount int         `json:"commentCount,omitempty"`
	UserName     string      `json:"username,omitempty"`
	Sort         *blog.Sort  `json:"sort,omitempty"`
	Label        *blog.Label `json:"label,omitempty"`
}

func (to *ArticleVO) Copy(from *blog.Article) {
	to.ID = from.ID
	to.UserId = from.UserId
	to.ArticleCover = from.ArticleCover
	to.ArticleTitle = from.ArticleTitle
	to.ArticleContent = from.ArticleContent
	to.CommentStatus = from.CommentStatus
	to.RecommendStatus = from.RecommendStatus
	to.Password = from.Password
	to.ViewStatus = from.ViewStatus
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
	to.UpdateBy = from.UpdateBy
	to.SortId = from.SortId
	to.LabelId = from.LabelId
}
