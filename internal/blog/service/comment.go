package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sync/atomic"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
	sDto "wan_go/pkg/common/dto"
	"wan_go/pkg/common/log"
	"wan_go/pkg/common/mail"
	"wan_go/pkg/utils"
	"wan_go/sdk/pkg"
	"wan_go/sdk/pkg/jwtauth/user"
	"wan_go/sdk/service"
)

type Comment struct {
	service.Service
}

func NewComment(c *gin.Context) *Comment {
	us := Comment{}
	us.Orm = pkg.Orm(c)
	us.Log = api.GetRequestLogger(c)
	return &us
}

func (s *Comment) Insert(c *gin.Context, d *dto.SaveCommentReq) error {

	data := blog.Comment{}
	d.CopyTo(&data)
	err := s.Orm.Debug().Model(&data).Create(&data).Error
	if err != nil {
		s.Log.Errorf("InsertReq error:%s", err)
		return err
	}

	//todo redis缓存

	if d.Type == blog_const.COMMENT_TYPE_ARTICLE.Code {
		article := blog.Article{}
		article.ID = d.Source
		if err = s.Orm.Debug().Select("user_id, article_title, comment_status").
			Find(&article).Error; err != nil {
			s.Log.Errorf("InsertReq error:%s", err)
			return err
		}
		if !article.CommentStatus {
			return errors.New("评论功能已关闭！")
		}

		mails, toName := GetMails(c, d, article.UserId, d.UserId)
		if len(mails) == 0 {
			return nil
		}
		SendMails(d, mails, toName, article.ArticleTitle, d.UserId, d.UserName)
	}

	return nil
}

func SendMails(d *dto.SaveCommentReq, mails []string, toName, articleTitle string, userId int32, userName string) {

	key := blog_const.COMMENT_IM_MAIL + mails[0]
	get, ok := cache.Get(key)
	if !ok {
		return
	}
	sendCount := get.(*atomic.Int32)

	var mailContent string
	comment0 := blog.Comment{ID: d.ParentCommentId}
	if err := db.Mysql().Select("comment_content").Find(&comment0); err == nil {
		mailContent = comment0.CommentContent
	}

	webName := cache.GetWebName()

	sendSuccess := func() {
		var i atomic.Int32
		i.Store(1)
		cache.SetExpire(key, &i, blog_const.TOKEN_EXPIRE*4)
	}
	mail.SendCommentMail(d, mails, articleTitle, userName, toName, webName, mailContent, sendCount, sendSuccess)
}

func GetMails(c *gin.Context, vo *dto.SaveCommentReq, articleUserId, currentUserId int32) (mails []string, toName string) {
	user := blog.User{}
	switch {
	case vo.ParentUserId > 0:
		user.ID = vo.ParentUserId
	case vo.Type == blog_const.COMMENT_TYPE_MESSAGE.Code:
		fallthrough
	case vo.Type == blog_const.COMMENT_TYPE_LOVE.Code:
		user.ID = int32(blog_const.ADMIN_USER_ID)
	case vo.Type == blog_const.COMMENT_TYPE_ARTICLE.Code:
		user.ID = articleUserId
	}

	us := NewUser(c)
	if err := us.GetUser(&user); err != nil {
		return
	}
	if user.ID > 0 && user.ID != currentUserId && len(user.Email) > 0 {
		toName = user.UserName
		mails = append(mails, user.Email)
	}
	return
}

func (s *Comment) Delete(d *dto.DelCommentReq, userId int32) error {
	data := blog.Comment{}
	err := s.Orm.Debug().Model(&data).
		Where("user_id=?", userId).Delete(&data, d.GetId()).Error
	if err != nil {
		s.Log.Errorf("Delete error: %s", err)
		return err
	}
	return nil
}

func commentCountKey(source int32, typ string) string {
	return fmt.Sprintf(blog_const.COMMENT_COUNT_CACHE + utils.Int32ToString(source) + "_" + typ)
}

func (s *Comment) CountComment(d *dto.CountCommentReq, count *int64) error {
	key := commentCountKey(d.Source, d.Type)

	if get, b := cache.Get(key); b {
		*count = get.(int64)
		return nil
	}
	if err := s.Orm.Debug().Model(&blog.Comment{}).
		Where("source = ? and type = ?", d.Source, d.Type).
		Count(count).Error; err != nil {
		log.NewWarn("GetCommentCount", err.Error())
		return err
	}

	cache.Set(key, *count)

	return nil
}

func (s *Comment) UserDeleteComment(d *dto.DelCommentReq, userId int32) error {

	comment := blog.Comment{}
	var err error
	if err = s.Orm.Debug().Select("source, type").Find(&comment, d.GetId()).Error; err != nil {
		s.Log.Errorf("UserDeleteComment error: %s", err)
		err = errors.New("数据不存在或无权删除该数据")
		return err
	}

	if comment.Type != blog_const.COMMENT_TYPE_ARTICLE.Code {
		err = errors.New("权限不足")
		return err
	}

	article := blog.Article{}
	if err = db.Mysql().Model(&blog.Article{}).Select("user_id").
		Where("id=?", comment.Source).First(&article).Error; err != nil {
		err = errors.New("文章不存在")
		return err
	}

	if userId != article.UserId {
		err = errors.New("无权删除他人评论")
		return err
	}

	err = s.Orm.Debug().Model(&comment).Delete(&comment, d.GetId()).Error
	if err != nil {
		s.Log.Errorf("Delete error: %s", err)
		return err
	}
	return nil
}

func (s *Comment) BossDeleteComment(d *dto.DelCommentReq) error {

	comment := blog.Comment{}
	err := s.Orm.Debug().Model(&comment).Delete(&comment, d.GetId()).Error
	if err != nil {
		s.Log.Errorf("Delete error: %s", err)
		return err
	}
	return nil
}

func (s *Comment) Page(context *gin.Context, d *dto.PageCommentReq, p *actions.DataPermission, page *vo.Page[vo.CommentVO]) error {

	if blog_const.COMMENT_TYPE_ARTICLE.Code == d.CommentType {
		article := blog.Article{ID: d.Source}
		if s.Orm.Debug().Select("comment_status").Find(&article).Error != nil || !article.CommentStatus {
			return errors.New("评论功能已关闭！")
		}
	}

	if d.FloorCommentId == 0 {
		//获取顶层评论
		var comments []blog.CommentExt
		err := s.Orm.Debug().Select("comment.*, user.user_name as UserName").Scopes(
			sDto.Paginate(d.Pagination),
			func(db *gorm.DB) *gorm.DB {
				return db.Joins("left join user on user.id = comment.user_id").
					Where("source=? and type=? and parent_comment_id=?", d.Source, d.CommentType, blog_const.FIRST_COMMENT)
			},
		).Find(&comments).
			Limit(-1).
			Offset(-1).
			Count(&page.Total).Error
		if err != nil {
			s.Log.Errorf("PageComment error: %s", err)
			return err
		}
		//获取各楼层评论
		commentVOs := make([]vo.CommentVO, 0)
		for _, comment := range comments {
			commentVO := buildCommentVO(&comment, context)
			var childComments []blog.CommentExt
			p2 := vo.Page[vo.CommentVO]{}
			p2.Current, p2.Size, p2.Column = 1, 5, "created_at"
			s.Orm.Debug().Scopes(
				sDto.Paginate(&p2.Pagination),
				func(db *gorm.DB) *gorm.DB {
					return db.Where("source=? and type=? and floor_comment_id=?", d.Source, d.CommentType, comment.ID)
				},
			).Find(&childComments).
				Limit(-1).
				Offset(-1).
				Count(&p2.Total)

			if len(childComments) > 0 {
				ccVOs := make([]vo.CommentVO, 0)
				for _, cc := range childComments {
					ccVOs = append(ccVOs, *buildCommentVO(&cc, context))
				}
				p2.Records = ccVOs
			}
			commentVO.ChildComments = p2
			commentVOs = append(commentVOs, *commentVO)
		}
		page.Records = commentVOs
	} else {
		//获取楼层评论
		var childComments []blog.CommentExt
		if err := s.Orm.Debug().Scopes(
			sDto.Paginate(d.Pagination),
			func(db *gorm.DB) *gorm.DB {
				return db.Where("source=? and type=? and floor_comment_id=?", d.Source, d.CommentType, d.FloorCommentId)
			},
		).Find(&childComments).Error; err != nil {
			s.Log.Errorf("PageComment error: %s", err)
			return err
		}

		if len(childComments) > 0 {
			ccVOs := make([]vo.CommentVO, 0)
			for _, cc := range childComments {
				cc2 := cc
				ccVOs = append(ccVOs, *buildCommentVO(&cc2, context))
			}
			page.Records = ccVOs
		}
	}

	return nil
}

func buildCommentVO(comment *blog.CommentExt, c *gin.Context) *vo.CommentVO {

	ret := vo.CommentVO{}
	ret.Copy(comment)

	//todo 缓存user基本信息 从缓存而不是从db获取
	var user blog.User
	//jwt.GetUser(c, &user)
	if user.ID > 0 {
		ret.Avatar = user.Avatar
		ret.UserName = user.UserName
	}

	if utils.IsEmpty(ret.UserName) {
		ret.UserName = utils.RandomName(ret.UserId)
	}

	if ret.ParentUserId > 0 {
		//todo
		var pu blog.User
		//jwt.GetUser(c, &pu)
		//pu := db_common.GetUser(ret.ParentUserId)
		if pu.ID > 0 {
			ret.ParentUserName = pu.UserName
		}
		if len(ret.ParentUserName) == 0 {
			ret.ParentUserName = utils.RandomName(ret.ParentUserId)
		}
	}

	return &ret
}

func (s *Comment) PageAdmin(c *gin.Context, d *dto.PageAdminCommentReq, p *actions.DataPermission, page *vo.Page[blog.Comment],
	isBoss bool) error {

	userId := user.GetUserId32(c)

	if !isBoss {
		as := NewArticle(c)
		userArticleIds := as.GetUserArticleIds(userId)
		if userArticleIds == nil {
			page.Records = []blog.Comment{}
			return nil
		}

		d.CommentType = blog_const.COMMENT_TYPE_ARTICLE.Code
		if d.Source == 0 {
			d.Sources = userArticleIds
		}
	}

	var data blog.Comment
	err := s.Orm.Debug().Scopes(
		sDto.MakeCondition(d.GetNeedSearch()),
		sDto.Paginate(d.Pagination),
		actions.Permission(data.TableName(), p),
	).Find(&page.Records).
		Limit(-1).
		Offset(-1).
		Count(&page.Total).Error

	if err != nil {
		s.Log.Errorf("PageUser error: %s", err)
		return err
	}

	return nil
}
