package dto

import (
	"time"
	"wan_go/pkg/common/db/mysql/blog"
	r "wan_go/pkg/common/response"
)

// 更改部分字段
type ChangeArticleReq struct {
	ArticleId       int32 `uri:"articleId"`
	ViewStatus      *bool `json:"viewStatus"`
	CommentStatus   *bool `json:"commentStatus"`
	RecommendStatus *bool `json:"recommendStatus"`
}

func (s *ChangeArticleReq) GetId() interface{} {
	return s.ArticleId
}

func (from *ChangeArticleReq) CopyTo(to *blog.Article) {
	if from.ArticleId != 0 {
		to.ID = from.ArticleId
	}
	if from.ViewStatus != nil {
		to.ViewStatus = *from.ViewStatus
	}
	if from.CommentStatus != nil {
		to.CommentStatus = *from.CommentStatus
	}
	if from.RecommendStatus != nil {
		to.RecommendStatus = *from.RecommendStatus
	}
}

type SaveArticleReq struct {
	ID              int32       `uri:"id"`
	UserId          int32       `json:"userId"`
	ArticleCover    string      `json:"articleCover"`
	ArticleTitle    string      `json:"articleTitle" vd:"@:len($)>0; msg:'用户名不能为空'"`
	ArticleContent  string      `json:"articleContent" vd:"@:len($)>0; msg:'用户名不能为空'"`
	ViewCount       int         `json:"viewCount"`
	LikeCount       int         `json:"likeCount"`
	CommentStatus   bool        `json:"commentStatus"`
	RecommendStatus bool        `json:"recommendStatus"`
	Password        string      `json:"password"`
	ViewStatus      bool        `json:"viewStatus"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updateTime"`
	UpdateBy        string      `json:"updateBy"`
	SortId          int32       `json:"sortId" vd:"$>0; msg:'文章分类不能为空'"`
	LabelId         int32       `json:"labelId" vd:"$>0; msg:'文章标签不能为空'"`
	CommentCount    int         `json:"commentCount"`
	UserName        string      `json:"username"`
	Sort            *blog.Sort  `json:"sort"`
	Label           *blog.Label `json:"label"`
}

func (from *SaveArticleReq) CopyTo(to *blog.Article) {
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

func (s *SaveArticleReq) GetId() interface{} {
	return s.ID
}

type DelArticleReq struct {
	Ids []int `json:"ids"`
}

func (s *DelArticleReq) GetId() interface{} {
	return s.Ids
}

type PageArticleReq struct {
	*r.Pagination   `json:",inline"`
	ArticleId       int32  `form:"articleId" search:"type:eq;column:id;table:article"`
	UserId          int32  `uri:"userId" search:"type:eq;column:user_id;table:article"`
	ArticleTitle    string `form:"articleTitle" search:"type:eq;column:article_title;table:article"`
	RecommendStatus *bool  `form:"recommendStatus" search:"type:eq;column:recommend_status;table:article"`
	LabelId         int32  `form:"labelId" search:"type:eq;column:label_id;table:article"`
	ViewStatus      *bool  `form:"viewStatus" search:"type:eq;column:view_status;table:article"`
}

func (s *PageArticleReq) GetNeedSearch() interface{} {
	return *s
}

type GetArticleReq struct {
	ID         int32  `form:"id"`
	ViewStatus *bool  `form:"viewStatus,omitempty"`
	Password   string `form:"password,omitempty"`
}
