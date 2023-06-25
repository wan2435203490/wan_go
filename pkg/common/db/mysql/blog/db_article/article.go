package db_article

import (
	"github.com/pkg/errors"
	"strings"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/db/mysql/blog/db_common"
	"wan_go/pkg/common/log"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
)

func UpdateVO(data *blogVO.ArticleVO) error {
	if !data.ViewStatus && utils.IsEmpty(data.Password) {
		return errors.New("请设置文章密码")
	}

	user := cache.GetUser()
	if user == nil {
		return errors.New("用户信息缓存失效，请重新登录")
	}

	updates := blog.Article{
		LabelId:         data.LabelId,
		SortId:          data.SortId,
		ArticleTitle:    data.ArticleTitle,
		UpdateBy:        user.UserName,
		ArticleContent:  data.ArticleContent,
		CommentStatus:   data.CommentStatus,
		RecommendStatus: data.RecommendStatus,
		ViewStatus:      data.ViewStatus,
	}

	if utils.IsNotEmpty(data.ArticleCover) {
		updates.ArticleCover = data.ArticleCover
	}

	if !data.ViewStatus && utils.IsNotEmpty(data.Password) {
		updates.Password = data.Password
	}

	err := db.Mysql().Where("id=? and user_id = ?", data.ID, user.ID).
		Updates(updates).Error

	db_common.CacheSort()

	return err
}

func InsertVO(data *blogVO.ArticleVO) error {

	if !data.ViewStatus && utils.IsEmpty(data.Password) {
		return errors.New("请设置文章密码")
	}

	article := blog.Article{}
	article.ArticleCover = data.ArticleCover
	if !data.ViewStatus {
		article.Password = data.Password
	}
	article.ViewStatus = data.ViewStatus
	article.CommentStatus = data.CommentStatus
	article.RecommendStatus = data.RecommendStatus
	article.ArticleTitle = data.ArticleTitle
	article.ArticleContent = data.ArticleContent
	article.SortId = data.SortId
	article.LabelId = data.LabelId
	article.UserId = data.UserId

	if err := Insert(&article); err != nil {
		return err
	}

	db_common.CacheSort()

	return nil
}

func Insert(data *blog.Article) error {
	return db.Mysql().Create(&data).Error
}

func DeleteByUserId(id int) error {
	err := db.Mysql().Where("user_id = ?", cache.GetUserId()).Delete(&blog.Article{ID: int32(id)}).Error
	db_common.CacheSort()
	return err
}

func Delete(data *blog.Article) error {
	return db.Mysql().Delete(&data).Error
}

func ListArticle(vo *blogVO.BaseRequestVO[*blogVO.ArticleVO]) {

	tx := db.Page(&vo.Pagination)

	if utils.IsNotEmpty(vo.SearchKey) {
		tx = tx.Where("article_title=?", vo.SearchKey)
	}
	if vo.RecommendStatus {
		tx = tx.Where("recommend_status=?", true)
	}

	if vo.LabelId > 0 {
		tx = tx.Where("label_id=?", vo.LabelId)
	} else if vo.SortId > 0 {
		tx = tx.Where("sort_id=?", vo.SortId)
	}

	var articles []*blog.Article
	if err := tx.
		Where("view_status = ?", true).
		Find(&articles).Error; err != nil {
		log.NewWarn("ListArticle", err.Error())
		return
	}

	var articlesVO []*blogVO.ArticleVO
	if len(articles) > 0 {
		for _, art := range articles {
			art.Password = ""
			if len(art.ArticleContent) > blog_const.SUMMARY {
				summary := art.ArticleContent[0:blog_const.SUMMARY]
				summary = strings.ReplaceAll(summary, "`", "")
				summary = strings.ReplaceAll(summary, "#", "")
				summary = strings.ReplaceAll(summary, ">", "")
				art.ArticleContent = summary
			}

			articleVO := buildArticleVO(art, false)
			articlesVO = append(articlesVO, articleVO)
		}
	}

	vo.Total = len(articlesVO)
	vo.Records = articlesVO

	return
}

func ExistArticleByUserId(id int32) bool {
	var count int64
	db.Mysql().Where("id=? and user_id=?", id, cache.GetUserId()).Count(&count)
	return count > 0
}

func GetArticleById(id int, flag bool, password string) *blogVO.ArticleVO {

	var vo blogVO.ArticleVO

	article := blog.Article{ID: int32(id)}
	tx := db.Mysql()
	if flag {
		tx.Where("view_status=?", true)
	} else {
		if utils.IsEmpty(password) {
			vo.Msg = "请输入文章密码！"
			return &vo
		}
		tx.Where("password=?", password)
	}
	if tx.Omit("password").Find(&article); tx.Error != nil || tx.RowsAffected == 0 {
		return nil
	}

	if err := db.Mysql().Model(&article).Updates("view_count=view_count+1").Error; err != nil {
		vo.Msg = err.Error()
		return &vo
	}

	//vo = *(buildArticleVO(&article, false))
	return buildArticleVO(&article, false)
}

func ListAdminArticle(vo *blogVO.BaseRequestVO[*blogVO.ArticleVO], isBoss bool) {
	tx := db.Page(&vo.Pagination).Omit("article_content, password")

	if isBoss {
		if vo.UserId > 0 {
			tx.Where("user_id=?", vo.UserId)
		}
	} else {
		tx.Where("user_id=?", cache.GetUserId())
	}

	if utils.IsNotEmpty(vo.SearchKey) {
		tx.Where("article_title=?", vo.SearchKey)
	}
	if vo.RecommendStatus {
		tx.Where("recommend_status=?", true)
	}
	if vo.LabelId > 0 {
		tx.Where("label_id=?", vo.LabelId)
	}
	if vo.SortId > 0 {
		tx.Where("sort_id=?", vo.SortId)
	}
	var articles []*blog.Article
	if err := tx.Debug().Model(blog.Article{}).Order("created_at desc").Find(&articles).Error; err != nil {
		vo.Msg = err.Error()
		return
	}

	if len(articles) > 0 {
		for _, article := range articles {
			articleVO := buildArticleVO(article, true)
			vo.Records = append(vo.Records, articleVO)
		}
	}
	return
}

func ChangeArticleStatus(id int, viewStatus, commentStatus, recommendStatus, viewStatusExist, commentStatusExist, recommendStatusExist bool) {
	tx := db.Mysql().Model(&blog.Article{ID: int32(id)}).Where("user_id=?", cache.GetUserId())

	updateColumns := make(map[string]any, 4)
	if viewStatusExist {
		updateColumns["view_status"] = viewStatus
	}
	if commentStatusExist {
		updateColumns["comment_status"] = commentStatus
	}
	if recommendStatusExist {
		updateColumns["recommend_status"] = recommendStatus
	}
	tx.UpdateColumns(updateColumns)
}

func GetArticleByIdForUser(id int) *blogVO.ArticleVO {
	var vo blogVO.ArticleVO
	article := blog.Article{ID: int32(id)}
	if err := db.Mysql().Where("user_id=?", cache.GetUserId()).Find(&article).Error; err != nil {
		vo.Msg = err.Error()
		return &vo
	}

	vo.Copy(&article)

	return &vo
}

func buildArticleVO(article *blog.Article, isAdmin bool) *blogVO.ArticleVO {

	vo := blogVO.ArticleVO{}
	vo.Copy(article)

	if !isAdmin && utils.IsEmpty(vo.ArticleCover) {
		//todo
		//vo.ArticleCover = randomcover
	}

	user := db_common.GetUser(vo.UserId)
	if user != nil && utils.IsNotEmpty(user.UserName) {
		vo.UserName = user.UserName
	} else if !isAdmin {
		//todo
		//vo.UserName = randomname
	}

	if vo.CommentStatus {
		vo.CommentCount = db_common.GetCommentCount(vo.ID, blog_const.COMMENT_TYPE_ARTICLE.Code)
	}

	get, b := cache.Get(blog_const.SORT_INFO)
	if !b || get == nil {
		return &vo
	}
	sortInfo := get.(*[]*blog.Sort)

	for _, si := range *sortInfo {
		if si.ID == vo.SortId {
			sort := blog.Sort{}
			sort.Copy(si)
			sort.Labels = nil
			vo.Sort = &sort

			if len(si.Labels) > 0 {
				for _, l := range si.Labels {
					if l.ID == vo.LabelId {
						label := blog.Label{}
						label.Copy(l)
						vo.Label = &label
						break
					}
				}
			}
			break
		}
	}

	return &vo
}
