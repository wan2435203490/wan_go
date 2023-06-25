package db_common

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/db/mysql/blog/db_user"
	"wan_go/pkg/common/log"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
)

func DB() *gorm.DB {
	return db.DB.MysqlDB
}

func GetUser(userId int32) *blog.User {

	if get, b := cache.Get(blog_const.USER_CACHE); b {
		return get.(*blog.User)
	}

	u := blog.User{ID: userId}
	err := db_user.Get(&u)
	if err != nil {
		log.NewWarn("GetUser", err.Error())
		return nil
	}

	return &u
}

func GetAdmire() []*blog.User {

	if get, b := cache.Get(blog_const.ADMIRE); b {
		return get.([]*blog.User)
	}

	var users []*blog.User

	if err := DB().
		Where("admire is not null").
		Select("user_name", "admire", "avatar").
		Find(&users).
		Error; err != nil {

		log.NewWarn("GetAdmire", err.Error())
		return nil
	}

	cache.Set(blog_const.ADMIRE, users)

	return users
}

func GetFamilyList() []*blogVO.FamilyVO {
	if get, b := cache.Get(blog_const.FAMILY_LIST); b {
		return get.([]*blogVO.FamilyVO)
	}

	var families []*blogVO.FamilyVO
	if err := DB().Model(&blog.Family{}).
		// todo test
		Where("status = ?", true).
		Find(&families).Error; err != nil {

		//if !errors.Is() todo filter no record error
		log.NewWarn("GetFamilyList", err.Error())
		return nil
	}

	cache.Set(blog_const.FAMILY_LIST, families)

	return families
}

func commentCountKey(source int32, typ string) string {
	return fmt.Sprintf(blog_const.COMMENT_COUNT_CACHE + utils.Int32ToString(source) + "_" + typ)
}

func GetCommentCount(source int32, typ string) int {

	key := commentCountKey(source, typ)

	if get, b := cache.Get(key); b {
		return get.(int)
	}

	var count int64
	if err := DB().Model(&blog.Comment{}).
		Where("source = ? and type = ?", source, typ).
		Count(&count).Error; err != nil {

		log.NewWarn("GetCommentCount", err.Error())
		return 0
	}

	cache.Set(key, int(count))

	return int(count)
}

func GetUserArticleIds(userId int) *[]int {

	key := blog_const.COMMENT_COUNT_CACHE + strconv.Itoa(userId)
	if get, b := cache.Get(key); b {
		return get.(*[]int)
	}

	var ret []int
	if err := DB().Model(&blog.Article{}).
		Where("userId = ?", userId).
		Select("id").
		Find(&ret).Error; err != nil {

		log.NewWarn("GetUserArticleIds", err.Error())
		return nil
	}

	cache.Set(key, &ret)

	return &ret
}

// CountSort CountLabel Maybe use one struct ?
type CountSort struct {
	SortId    int32 `gorm:"column:sort_id"`
	SortCount int32
}

type CountLabel struct {
	LabelId    int32 `gorm:"column:label_id"`
	LabelCount int32
}

func CacheSort() {
	info := GetSortInfo()
	if len(*info) > 0 {
		cache.Set(blog_const.SORT_INFO, info)
	} else {
		cache.Delete(blog_const.SORT_INFO)
	}
}

func GetSortInfo() *[]*blog.Sort {

	var sorts []*blog.Sort
	if err := DB().Find(&sorts).Error; err != nil {
		log.NewWarn("GetSortInfo", err.Error())
		return nil
	}

	var countSorts []*CountSort
	if err := DB().Model(&blog.Article{}).
		Select("sort_id, count(1) as SortCount").
		Where("view_status = ?", blog_const.STATUS_ENABLE.Code).
		Group("sort_id").
		Find(&countSorts).
		Error; err != nil {

		log.NewWarn("GetSortInfo-countSorts", err.Error())
		return nil
	}

	var countLabels []*CountLabel
	if err := DB().Model(&blog.Article{}).
		Select("label_id, count(1) as LabelCount").
		Where("view_status = ?", blog_const.STATUS_ENABLE.Code).
		Group("label_id").
		Find(&countLabels).
		Error; err != nil {

		log.NewWarn("GetSortInfo-countLabels", err.Error())
		return nil
	}

	var labels []*blog.Label
	//直接搜的全量label 数据量大考虑增加过滤sort_id
	if err := DB().Find(&labels).Error; err != nil {
		//todo no record validate
		log.NewWarn("GetSortInfo-Labels", err.Error())
		return nil
	}

	if sorts == nil {
		return nil
	}

	for _, sort := range sorts {

		for _, countSort := range countSorts {
			if sort.ID == countSort.SortId {
				sort.CountOfSort = countSort.SortCount
			}
		}

		var sortLabels []*blog.Label
		for _, label := range labels {
			if label.SortId == sort.ID {
				sortLabels = append(sortLabels, label)

				for _, countLabel := range countLabels {
					if countLabel.LabelId == label.ID {
						label.CountOfLabel = countLabel.LabelCount
					}
				}
			}
		}
		sort.Labels = sortLabels
	}

	return &sorts
}
