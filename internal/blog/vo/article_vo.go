package vo

import (
	"time"
	"wan_go/pkg/common/db/mysql/blog"
	r "wan_go/pkg/common/response"
)

type ArticleVO struct {
	r.CodeMsg
	ID     int32 `json:"id"`
	UserId int32 `json:"userId"`
	// 查询为空时，随机选择
	ArticleCover    string    `json:"articleCover"`
	ArticleTitle    string    `json:"articleTitle" vd:"@:len($)>0; msg:'用户名不能为空'"`
	ArticleContent  string    `json:"articleContent" vd:"@:len($)>0; msg:'用户名不能为空'"`
	ViewCount       int       `json:"viewCount"`
	LikeCount       int       `json:"likeCount"`
	CommentStatus   bool      `json:"commentStatus"`
	RecommendStatus bool      `json:"recommendStatus"`
	Password        string    `json:"password"`
	ViewStatus      bool      `json:"viewStatus"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updateTime"`
	UpdateBy        string    `json:"updateBy"`
	SortId          int32     `json:"sortId" vd:"$>0; msg:'文章分类不能为空'"`
	LabelId         int32     `json:"labelId" vd:"$>0; msg:'文章标签不能为空'"`
	// 需要查询封装
	CommentCount int64       `json:"commentCount"`
	UserName     string      `json:"username"`
	Sort         *blog.Sort  `json:"sort"`
	Label        *blog.Label `json:"label"`
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
