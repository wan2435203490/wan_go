package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"wan_go/internal/blog/service/dto"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/actions"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"
	rocksCache "wan_go/pkg/common/db/rocks_cache"
	sDto "wan_go/pkg/common/dto"
	"wan_go/pkg/utils"
	"wan_go/sdk/pkg"
	"wan_go/sdk/service"
)

type Article struct {
	service.Service
}

func NewArticle(c *gin.Context) *Article {
	us := Article{}
	us.Orm = pkg.Orm(c)
	us.Log = api.GetRequestLogger(c)
	return &us
}

func (s *Article) Update(d *dto.SaveArticleReq) error {

	var model blog.Article
	var count int64
	if err := s.Orm.Debug().Model(&model).Where("id<>? and article_title=?", d.ID, d.ArticleTitle).Count(&count).Error; err != nil {
		s.Log.Errorf("Update Count error:%s", err)
		return err
	}
	if count > 0 {
		return errors.New("文章名称不能重复！")
	}

	d.CopyTo(&model)
	tx := s.Orm.Debug().Model(&model).Updates(&d)
	if err := tx.Error; err != nil {
		s.Log.Errorf("Update error:%s", err)
		return err
	}
	if tx.RowsAffected == 0 {
		return errors.New("数据不存在或无权更新该数据")
	}
	return nil
}

func (s *Article) Insert(d *dto.SaveArticleReq) error {

	var model blog.Article
	var count int64
	if err := s.Orm.Debug().Model(&model).Where("article_title=?", d.ArticleTitle).Count(&count).Error; err != nil {
		s.Log.Errorf("InsertReq Count error:%s", err)
		return err
	}
	if count > 0 {
		return errors.New("文章名称不能重复！")
	}

	data := blog.Article{}
	d.CopyTo(&data)
	err := s.Orm.Debug().Model(&model).Create(&data).Error
	if err != nil {
		s.Log.Errorf("InsertReq error:%s", err)
		return err
	}
	return nil
}

func (s *Article) Delete(d *dto.DelArticleReq) error {
	var err error
	var data blog.Article

	tx := s.Orm.Debug().Model(&data).Delete(&data, d.GetId())

	if err = tx.Error; err != nil {
		err = tx.Error
		s.Log.Errorf("Delete error: %s", err)
		return err
	}
	if tx.RowsAffected == 0 {
		err = errors.New("数据不存在或无权删除该数据")
		return err
	}
	return nil
}

func (s *Article) GetArticle(c *gin.Context, d *dto.GetArticleReq, data *vo.ArticleVO, userName string) error {
	article := blog.Article{ID: d.ID}
	tx := s.Orm.Debug()
	if d.ViewStatus != nil {
		tx.Where("view_status=?", d.ViewStatus)
	} else {
		if utils.IsEmpty(d.Password) {
			return errors.New("请输入文章密码")
		}
		tx.Where("password=?", d.Password)
	}
	if tx.Omit("password").Find(&article); tx.Error != nil {
		return errors.New("文章不存在")
	}

	s.Orm.Debug().Model(&article).Updates("view_count=view_count+1")

	*data = *buildArticleVO(c, &article, false, userName)

	return nil
}

func (s *Article) Page(d *dto.PageArticleReq, p *actions.DataPermission, page *vo.Page[blog.Article]) error {

	var data blog.Article
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

func (s *Article) ListArticle(c *gin.Context, d *dto.PageArticleReq, p *actions.DataPermission, page *vo.Page[vo.ArticleVO],
	userName string) error {

	viewStatus := true
	d.ViewStatus = &viewStatus

	var data blog.Article
	var articles []blog.Article
	err := s.Orm.Debug().Scopes(
		sDto.MakeCondition(d.GetNeedSearch()),
		sDto.Paginate(d.Pagination),
		//fun,
		actions.Permission(data.TableName(), p),
	).Find(&articles).
		Limit(-1).
		Offset(-1).
		Count(&page.Total).Error
	if err != nil {
		s.Log.Errorf("ListArticle error: %s", err)
		return err
	}

	var articlesVO []vo.ArticleVO
	if len(articles) > 0 {
		for _, v := range articles {
			art := v
			art.Password = ""
			if len(art.ArticleContent) > blog_const.SUMMARY {
				summary := art.ArticleContent[0:blog_const.SUMMARY]
				summary = strings.ReplaceAll(summary, "`", "")
				summary = strings.ReplaceAll(summary, "#", "")
				summary = strings.ReplaceAll(summary, ">", "")
				art.ArticleContent = summary
			}

			articleVO := buildArticleVO(c, &art, false, userName)
			articlesVO = append(articlesVO, *articleVO)
		}
	}
	page.Records = articlesVO

	return nil
}

func (s *Article) ListAdminArticle(c *gin.Context, d *dto.PageArticleReq, p *actions.DataPermission, page *vo.Page[vo.ArticleVO],
	isBoss bool, userId int32, userName string) error {

	if !isBoss {
		d.UserId = userId
	}

	var data blog.Article
	var articles []blog.Article
	err := s.Orm.Debug().Omit("article_content, password").Scopes(
		sDto.MakeCondition(d.GetNeedSearch()),
		sDto.Paginate(d.Pagination),
		//fun,
		actions.Permission(data.TableName(), p),
	).Find(&articles).
		Limit(-1).
		Offset(-1).
		Count(&page.Total).Error
	if err != nil {
		s.Log.Errorf("ListAdminArticle error: %s", err)
		return err
	}

	if len(articles) > 0 {
		for _, article := range articles {
			art := article
			articleVO := buildArticleVO(c, &art, true, userName)
			page.Records = append(page.Records, *articleVO)
		}
	}

	return nil
}

func buildArticleVO(c *gin.Context, article *blog.Article, isAdmin bool, userName string) *vo.ArticleVO {

	data := vo.ArticleVO{}
	data.Copy(article)

	if !isAdmin && utils.IsEmpty(data.ArticleCover) {
		data.ArticleCover = utils.RandomCover(0)
	}

	if utils.IsNotEmpty(userName) {
		data.UserName = userName
	} else if !isAdmin {
		data.UserName = utils.RandomName(0)
	}

	sortInfo, err := rocksCache.GetSortInfo()
	if err != nil || sortInfo == nil {
		return &data
	}

	for _, si := range sortInfo {
		if si.ID == data.SortId {
			sort := blog.Sort{}
			sort.Copy(si)
			sort.Labels = nil
			data.Sort = &sort

			if si.Labels != nil {
				for _, l := range *si.Labels {
					if l.ID == data.LabelId {
						label := blog.Label{}
						label.Copy(l)
						data.Label = &label
						break
					}
				}
			}
			break
		}
	}

	if data.CommentStatus {
		cs := NewComment(c)
		req := dto.CountCommentReq{Source: data.ID, Type: blog_const.COMMENT_TYPE_ARTICLE.Code}
		if err := cs.CountComment(&req, &data.CommentCount); err != nil {
			cs.Log.Errorf("CountComment error:%s", err)
			//maybe mustn't return error?
		}
	}

	return &data
}

func (s *Article) ChangeArticleStatus(d *dto.ChangeArticleReq, userId int32) error {
	tx := s.Orm.Debug().Model(&blog.Article{ID: d.ArticleId}).Where("user_id=?", userId)

	updateColumns := make(map[string]any, 4)
	if d.ViewStatus != nil {
		updateColumns["view_status"] = *d.ViewStatus
	}
	if d.CommentStatus != nil {
		updateColumns["comment_status"] = *d.CommentStatus
	}
	if d.RecommendStatus != nil {
		updateColumns["recommend_status"] = *d.RecommendStatus
	}
	tx = tx.UpdateColumns(updateColumns)

	if err := tx.Error; err != nil {
		s.Log.Errorf("Update error:%s", err)
		return err
	}
	if tx.RowsAffected == 0 {
		return errors.New("数据不存在或无权更新该数据")
	}
	return nil
}

func (s *Article) GetArticleByIdForUser(articleId, userId int32, data *vo.ArticleVO) error {

	article := blog.Article{ID: articleId}
	if err := s.Orm.Debug().Where("user_id=?", userId).Find(&article).Error; err != nil {
		return err
	}

	data.Copy(&article)

	return nil
}

func (s *Article) ExistArticleByUserId(articleId, userId int32, exist *bool) error {
	var count int64
	var model blog.Article
	if err := s.Orm.Debug().Model(&model).
		Where("id=? and user_id=?", articleId, userId).
		Count(&count).Error; err != nil {
		s.Log.Errorf("ExistArticleByUserId Count error:%s", err)
		return err
	}
	*exist = count > 0

	return nil
}

func (s *Article) GetUserArticleIds(userId int32) *[]int {

	key := blog_const.COMMENT_COUNT_CACHE + utils.Int32ToString(userId)
	if get, b := cache.Get(key); b {
		return get.(*[]int)
	}

	var ret []int
	if err := s.Orm.Debug().Model(&blog.Article{}).
		Where("userId = ?", userId).
		Select("id").
		Find(&ret).Error; err != nil {

		s.Log.Errorf("GetUserArticleIds error: %s", err)
		return nil
	}

	cache.Set(key, &ret)

	return &ret
}
