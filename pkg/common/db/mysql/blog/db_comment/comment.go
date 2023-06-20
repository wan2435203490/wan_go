package db_comment

import (
	"sync/atomic"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/db/mysql/blog/db_common"
	"wan_go/pkg/common/mail"
	r "wan_go/pkg/common/response"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
)

func ApiErr(msg string) *r.CodeMsg {
	return &r.CodeMsg{Msg: msg}
}

func SaveComment(vo *blogVO.CommentVO) *r.CodeMsg {
	if !blog_const.ExistsCommentType(vo.Type) {
		return ApiErr("评论来源类型不存在！")
	}

	article := blog.Article{}
	if vo.Type == blog_const.COMMENT_TYPE_ARTICLE.Code {
		article.ID = vo.Source
		if err := db.Mysql().Select("user_id, article_title, comment_status").
			Find(&article).Error; err != nil {
			return ApiErr(err.Error())
		}
		if !article.CommentStatus {
			return ApiErr("评论功能已关闭！")
		}
	}

	comment := blog.Comment{
		Source:          vo.Source,
		Type:            vo.Type,
		CommentContent:  vo.CommentContent,
		ParentCommentId: vo.ParentCommentId,
		FloorCommentId:  vo.FloorCommentId,
		ParentUserId:    vo.ParentUserId,
		UserId:          vo.UserId,
	}
	if utils.IsNotEmpty(vo.CommentInfo) {
		comment.CommentInfo = vo.CommentInfo
	}
	if err := Insert(&comment); err != nil {
		return ApiErr(err.Error())
	}

	mails, toName := GetMails(vo, article.UserId)
	if len(mails) == 0 {
		return nil
	}
	SendMails(vo, mails, toName, article.ArticleTitle)

	return nil
}

func SendMails(vo *blogVO.CommentVO, mails []string, toName, articleTitle string) {

	key := blog_const.COMMENT_IM_MAIL + mails[0]
	get, ok := cache.Get(key)
	if !ok {
		return
	}
	sendCount := get.(*atomic.Int32)

	var mailContent string
	comment0 := blog.Comment{ID: vo.ParentCommentId}
	if err := db.Mysql().Select("comment_content").Find(&comment0); err == nil {
		mailContent = comment0.CommentContent
	}

	webName := cache.GetWebName()
	fromName := cache.GetUserName()

	sendSuccess := func() {
		var i atomic.Int32
		i.Store(1)
		cache.SetExpire(key, &i, blog_const.TOKEN_EXPIRE*4)
	}
	mail.SendCommentMail(vo, mails, articleTitle, fromName, toName, webName, mailContent, sendCount, sendSuccess)
}

func GetMails(vo *blogVO.CommentVO, articleUserId int32) (mails []string, toName string) {
	//mails := make([]string, 0)
	var user *blog.User
	switch {
	case vo.ParentUserId > 0:
		user = db_common.GetUser(vo.ParentUserId)
	case vo.Type == blog_const.COMMENT_TYPE_MESSAGE.Code:
		fallthrough
	case vo.Type == blog_const.COMMENT_TYPE_LOVE.Code:
		user = cache.GetAdminUser()
	case vo.Type == blog_const.COMMENT_TYPE_ARTICLE.Code:
		user = db_common.GetUser(articleUserId)
	}

	return ResolveMails(user)
}

func ResolveMails(user *blog.User) (mails []string, toName string) {
	if cache.CanSendEmail(user) {
		toName = user.UserName
		mails = append(mails, user.Email)
	}
	return
}

func Insert(comment *blog.Comment) error {
	return db.Mysql().Create(comment).Error
}

func UserDeleteComment(id int) *r.CodeMsg {
	comment := blog.Comment{ID: int32(id)}
	if err := db.Mysql().Select("source, type").Find(&comment).Error; err != nil {
		return ApiErr(err.Error())
	}

	if comment.Type != blog_const.COMMENT_TYPE_ARTICLE.Code {
		return ApiErr("权限不足！")
	}

	var dest blog.Comment
	if err := db.Mysql().Model(&blog.Article{}).Select("user_id").
		Where("id=?", comment.Source).First(&dest).Error; err != nil {
		return ApiErr(err.Error())
	}

	if int32(cache.GetUserId()) != dest.UserId {
		return ApiErr("权限不足！")
	}

	Delete(id)

	return nil
}

func DeleteByUserId(id int) {
	db.Mysql().Where("user_id=?", cache.GetUserId()).Delete(&blog.Comment{ID: int32(id)})
}

func Delete(id int) {
	db.Mysql().Delete(&blog.Comment{ID: int32(id)})
}

func ListComment(vo *blogVO.BaseRequestVO[*blogVO.CommentVO]) {
	if vo.Source < 0 || len(vo.CommentType) == 0 {
		vo.CodeMsg = r.PARAMETER_ERROR
		return
	}
	if blog_const.COMMENT_TYPE_ARTICLE.Code == vo.CommentType {
		article := blog.Article{ID: vo.Source}
		if db.Mysql().Select("comment_status").Find(&article).Error != nil || !article.CommentStatus {
			vo.Msg = "评论功能已关闭！"
			return
		}
	}

	if vo.FloorCommentId == 0 {
		var comments []*blog.Comment
		if db.Page(&vo.Pagination).Where("source=? and type=? and parent_comment_id=?", vo.Source, vo.CommentType, blog_const.FIRST_COMMENT).
			Order("CreatedAt").
			Find(&comments).Error != nil || len(comments) == 0 {
			return
		}
		commentVOs := make([]*blogVO.CommentVO, 0)
		for _, c := range comments {
			commentVO := buildCommentVO(c)
			var childComments []*blog.Comment
			vo2 := blogVO.BaseRequestVO[*blogVO.CommentVO]{}
			vo2.Current, vo2.Size = 1, 5
			db.Page(&vo2.Pagination).
				Where("source=? and type=? and floor_comment_id=?", vo.Source, vo.CommentType, c.ID).
				Order("CreatedAt").
				Find(&childComments)

			if len(childComments) > 0 {
				ccVOs := make([]*blogVO.CommentVO, 0)
				for _, cc := range childComments {
					ccVOs = append(ccVOs, buildCommentVO(cc))
				}
				vo2.Records = ccVOs
			}
			//todo child 实现Pagination[T] 将Records放到Pagination里
			//commentVO.ChildComments = vo2
			commentVOs = append(commentVOs, commentVO)
		}
		vo.SetRecords(&commentVOs)

	} else {
		var childComments []*blog.Comment
		if db.Page(&vo.Pagination).
			Where("source=? and type=? and floor_comment_id=?", vo.Source, vo.CommentType, vo.FloorCommentId).
			Order("CreatedAt").
			Find(&childComments).Error != nil {
			return
		}

		if len(childComments) > 0 {
			ccVOs := make([]*blogVO.CommentVO, 0)
			for _, cc := range childComments {
				ccVOs = append(ccVOs, buildCommentVO(cc))
			}
			vo.Records = ccVOs
		}
	}
}

func ListAdminComment(vo *blogVO.BaseRequestVO[*blog.Comment], isBoss bool) {

	tx := db.Page(&vo.Pagination)

	var comments []*blog.Comment
	if isBoss {
		if vo.Source > 0 {
			tx.Where("source=?", vo.Source)
		}
		if utils.IsNotEmpty(vo.CommentType) {
			tx.Where("type=?", vo.CommentType)
		}
	} else {
		userArticleIds := db_common.GetUserArticleIds(cache.GetUserId())
		if userArticleIds == nil {
			vo.Records = []*blog.Comment{}
			return
		}
		if vo.Source > 0 {
			tx.Where("source=? and type=?", vo.Source, blog_const.COMMENT_TYPE_ARTICLE.Code)
		} else {
			tx.Where("type=? and source in ?", blog_const.COMMENT_TYPE_ARTICLE.Code, userArticleIds)
		}
	}

	tx.Order("CreatedAt DESC").Find(&comments)

	vo.SetRecords(&comments)
}

func buildCommentVO(c *blog.Comment) *blogVO.CommentVO {

	vo := blogVO.CommentVO{}
	vo.Copy(c)

	user := db_common.GetUser(vo.UserId)
	if user != nil {
		vo.Avatar = user.Avatar
		vo.UserName = user.UserName
	}

	if len(vo.UserName) == 0 {
		vo.UserName = utils.RandomName(vo.UserId)
	}

	if vo.ParentUserId > 0 {
		u := db_common.GetUser(vo.ParentUserId)
		if u != nil {
			vo.ParentUsername = u.UserName
		}
		if len(vo.ParentUsername) == 0 {
			vo.UserName = utils.RandomName(vo.ParentUserId)
		}
	}

	return &vo
}
