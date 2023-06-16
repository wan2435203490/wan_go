package db_resource_path

import (
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/log"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
)

func Get(data *blog.ResourcePath) error {
	return db.Mysql().Find(&data).Error
}

func List() ([]*blog.ResourcePath, error) {

	var datas []*blog.ResourcePath

	if err := db.Mysql().Find(&datas).Error; err != nil {
		return nil, err
	}

	return datas, nil
}

func Update(data *blog.ResourcePath) error {
	return db.Mysql().Updates(&data).Error
}

func Insert(data *blog.ResourcePath) error {
	return db.Mysql().Create(&data).Error
}

func Delete(data *blog.ResourcePath) error {
	return db.Mysql().Delete(&data).Error
}

func ListByResourceTypeAndClassify(requestVO *blogVO.BaseRequestVO[*blogVO.ResourcePathVO]) {

	//if !isAdmin {
	//	status = true
	//}

	var resourcePaths []*blog.ResourcePath

	//if err := db.Mysql().
	//	Where("type = ? and classify = ? and status = ?", resourceType, classify, status).
	//	Order(pagination.Order()).
	//	Limit(pagination.Size).
	//	Offset(pagination.Current).
	//	Find(&resourcePaths).Error; err != nil {
	//
	//	log.NewWarn("ListByResourceTypeAndClassify", err.Error())
	//	return nil
	//}

	if err := db.Page(&requestVO.Pagination).
		Where("type = ? and classify = ? and status = ?", requestVO.ResourceType, requestVO.Classify, requestVO.Status).
		Find(&resourcePaths).Error; err != nil {
		log.NewWarn("ListByResourceTypeAndClassify", err.Error())
		return
	}

	requestVO.Total = len(resourcePaths)

	if resourcePaths != nil {
		for _, path := range resourcePaths {
			p := blogVO.ResourcePathVO{}
			p.Copy(path)
			requestVO.Records = append(requestVO.Records, &p)
		}
	}

	return
}

// todo test: maybe can't find
func ListFunny() (ret []map[string]any) {
	if err := db.Mysql().Model(&blog.ResourcePath{}).
		Select("Classify, count(*) as Count").
		Where("status = ? and type = ?", true, blog_const.RESOURCE_PATH_TYPE_FUNNY).
		Group("classify").
		Find(&ret).Error; err != nil {
		return
	}

	return
}

func ListCollect() (classifyMap map[string][]blogVO.ResourcePathVO) {

	var paths []*blog.ResourcePath
	if err := db.Mysql().
		Where("status = ? and type = ?", true, blog_const.RESOURCE_PATH_TYPE_FAVORITES).
		//Group("classify").
		Order("classify, title").
		Find(&paths).Error; err != nil {
		return
	}

	//var classifyMap map[string][]blogVO.ResourcePathVO

	for _, path := range paths {
		//if _, ok := classifyMap[path.Classify]; !ok {
		var pathVO blogVO.ResourcePathVO
		pathVO.Copy(path)
		classifyMap[path.Classify] = append(classifyMap[path.Classify], pathVO)

		//}
	}

	return
}

func ListAdminLovePhoto(adminId int32) (ret []map[string]any) {
	//todo why it is []map?
	if err := db.Mysql().Model(&blog.ResourcePath{}).
		Select("classify, count(1) as Count").
		Where("status = ? and type = ? and remark = ?",
			true, blog_const.RESOURCE_PATH_TYPE_LOVE_PHOTO, utils.Int32ToString(adminId)).
		Group("classify").
		Find(&ret).Error; err != nil {
		return
	}

	return
}
