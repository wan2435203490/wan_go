package mysql

import (
	"errors"
	"gorm.io/gorm"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
)

func Orm() *gorm.DB {
	return db.DB.MysqlDB.Debug()
}

func GetAdmire() (users []blog.User, err error) {
	if err = Orm().
		Where("admire > 0").
		Select("user_name", "admire", "avatar").
		Find(&users).
		Error; err != nil {
		return nil, err
	}

	return users, nil
}

type CountSort struct {
	//gorm or json
	SortId    int32 `json:"sort_id,omitempty"`
	SortCount int32 `json:"sort_count,omitempty"`
}

type CountLabel struct {
	//gorm or json
	LabelId    int32 `gorm:"label_id,omitempty"`
	LabelCount int32 `gorm:"label_count,omitempty"`
}

func GetSortInfo() (sorts []*blog.Sort, err error) {

	if err := Orm().Model(&blog.Sort{}).Find(&sorts).Error; err != nil {
		return nil, err
	}

	if sorts == nil {
		return nil, errors.New("sortInfo is empty")
	}

	var countSorts []CountSort
	if err := Orm().Model(&blog.Article{}).
		Select("sort_id, count(1) as sort_count").
		Where("view_status = ?", blog_const.STATUS_ENABLE.Code).
		Group("sort_id").
		Find(&countSorts).
		Error; err != nil {
		//s.Log.Errorf("GetSortInfo countSorts error: %s", err)
		return nil, err
	}

	var countLabels []CountLabel
	if err := Orm().Model(&blog.Article{}).
		Select("label_id, count(1) as label_count").
		Where("view_status = ?", blog_const.STATUS_ENABLE.Code).
		Group("label_id").
		Find(&countLabels).
		Error; err != nil {
		//s.Log.Errorf("GetSortInfo countLabels error: %s", err)
		return nil, err
	}

	var labels []*blog.Label
	//直接搜的全量label 数据量大考虑增加过滤sort_id
	if err := Orm().Find(&labels).Error; err != nil {
		//s.Log.Errorf("GetSortInfo Labels error: %s", err)
		return nil, err
	}

	//或者传sorts *[]*Sort
	for _, sort := range sorts {

		for _, countSort := range countSorts {
			if sort.ID == countSort.SortId {
				sort.CountOfSort = countSort.SortCount
			}
		}

		var sortLabels []*blog.Label
		for _, label := range labels {
			sortLabels = append(sortLabels, label)

			for _, countLabel := range countLabels {
				if countLabel.LabelId == label.ID {
					label.CountOfLabel = countLabel.LabelCount
				}
			}
		}
		sort.Labels = &sortLabels
	}

	return sorts, nil
}

func GetWebInfo() ([]blog.WebInfo, error) {
	var webInfos []blog.WebInfo
	if err := Orm().Model(&blog.WebInfo{}).Find(&webInfos).Error; err != nil {
		return nil, err
	}
	return webInfos, nil
}

func GetByUserType(userType int8) (*blog.User, error) {
	user := blog.User{}
	err := Orm().Where("user_type=?", userType).Find(&user).Error
	return &user, err
}

func GetFamilyByUserId(userId int32) (*blog.Family, error) {
	family := blog.Family{}
	if err := Orm().Where("user_id=?", userId).First(&family).Error; err != nil {
		//s.Log.Errorf("GetByUserId error:%s", err)
		return nil, err
	}
	return &family, nil
}
