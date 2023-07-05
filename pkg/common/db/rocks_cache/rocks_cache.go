package rocksCache

//see https://github.com/dtm-labs/rockscache
//document
//every delete is only delete from cache

// use Fetch to fetch data
// 1. the first parameter is the key of the data
// 2. the second parameter is the data expiration time
// 3. the third parameter is the data fetch function which is called when the cache does not exist
//v, err := rc.Fetch("key1", 300 * time.Second, func()(string, error) {
//	// fetch data from database or other sources
//	return "value1", nil
//})

//Delete the cache
//rc.TagAsDeleted(key)

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"time"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/log"
	"wan_go/pkg/utils"
)

func DelKeys() {
	fmt.Println("init rocks cache to del old keys")
	fName := utils.GetSelfFuncName()
	for _, key := range []string{blog_const.ADMIN, blog_const.WEB_INFO, blog_const.SORT_INFO, blog_const.ADMIN_FAMILY,
		blog_const.ADMIRE} {
		var cursor uint64
		var n int
		for {
			var keys []string
			var err error
			keys, cursor, err = db.DB.RDB.Scan(context.Background(), cursor, key+"*", 3000).Result()
			if err != nil {
				panic(err.Error())
			}
			n += len(keys)
			// for each for redis cluster
			for _, key := range keys {
				if err = db.DB.RDB.Del(context.Background(), key).Err(); err != nil {
					log.NewError("", fName, key, err.Error())
					err = db.DB.RDB.Del(context.Background(), key).Err()
					if err != nil {
						panic(err.Error())
					}
				}
			}
			if cursor == 0 {
				break
			}
		}
	}
}

func GetWebInfo() ([]blog.WebInfo, error) {
	//网站信息 基本不会变
	webInfoStr, err := db.DB.Rc.Fetch(blog_const.WEB_INFO, time.Hour*30*24, func() (string, error) {
		info, err := mysql.GetWebInfo()
		if err != nil {
			return "", utils.Wrap(err)
		}
		str, err := sonic.MarshalString(info)
		if err != nil {
			return "", utils.Wrap(err)
		}
		return str, nil
	})
	if err != nil {
		return nil, utils.Wrap(err)
	}
	var webInfo []blog.WebInfo
	err = sonic.UnmarshalString(webInfoStr, &webInfo)
	return webInfo, utils.Wrap(err)
}

func DeleteWebInfo() error {
	err := db.DB.Rc.TagAsDeleted(blog_const.WEB_INFO)
	return utils.Wrap(err)
}

func GetSortInfo() ([]*blog.Sort, error) {
	sortInfoStr, err := db.DB.Rc.Fetch(blog_const.SORT_INFO, time.Hour*24, func() (string, error) {
		info, err := mysql.GetSortInfo()
		if err != nil {
			return "", utils.Wrap(err)
		}
		str, err := sonic.MarshalString(info)
		if err != nil {
			return "", utils.Wrap(err)
		}
		return str, nil
	})
	if err != nil {
		return nil, utils.Wrap(err)
	}
	var sortInfo []*blog.Sort
	err = sonic.UnmarshalString(sortInfoStr, &sortInfo)
	return sortInfo, utils.Wrap(err)
}

func DeleteSortInfo() error {
	err := db.DB.Rc.TagAsDeleted(blog_const.SORT_INFO)
	return utils.Wrap(err)
}

func GetAdminUser() (*blog.User, error) {
	userInfo, err := db.DB.Rc.Fetch(blog_const.ADMIN, time.Hour*24*30, func() (string, error) {
		user, err := mysql.GetByUserType(blog_const.UserRoleAdmin.Code)
		if err != nil {
			return "", utils.Wrap(err)
		}
		str, err := sonic.MarshalString(&user)
		if err != nil {
			return "", utils.Wrap(err)
		}
		return str, nil
	})
	if err != nil {
		return nil, utils.Wrap(err)
	}

	var user blog.User
	err = sonic.UnmarshalString(userInfo, &user)
	return &user, err
}

func DeleteAdminUser() error {
	err := db.DB.Rc.TagAsDeleted(blog_const.ADMIN)
	return utils.Wrap(err)
}

func GetAdminFamily(adminId int32) (*blog.Family, error) {
	str, err := db.DB.Rc.Fetch(blog_const.ADMIN_FAMILY, time.Hour*24, func() (string, error) {
		info, err := mysql.GetFamilyByUserId(adminId)
		if err != nil {
			return "", utils.Wrap(err)
		}
		str, err := sonic.MarshalString(info)
		if err != nil {
			return "", utils.Wrap(err)
		}
		return str, nil
	})
	if err != nil {
		return nil, utils.Wrap(err)
	}
	var family blog.Family
	err = sonic.UnmarshalString(str, &family)
	return &family, utils.Wrap(err)
}

func DeleteAdminFamily() error {
	err := db.DB.Rc.TagAsDeleted(blog_const.ADMIN_FAMILY)
	return utils.Wrap(err)
}

func GetAdmire() (*[]blog.User, error) {
	str, err := db.DB.Rc.Fetch(blog_const.ADMIRE, time.Hour*24, func() (string, error) {
		users, err := mysql.GetAdmire()
		if err != nil {
			return "", utils.Wrap(err)
		}
		str, err := sonic.MarshalString(users)
		if err != nil {
			return "", utils.Wrap(err)
		}
		return str, nil
	})
	if err != nil {
		return nil, utils.Wrap(err)
	}
	var users []blog.User
	err = sonic.UnmarshalString(str, &users)
	return &users, utils.Wrap(err)
}

func DeleteAdmire() error {
	err := db.DB.Rc.TagAsDeleted(blog_const.ADMIRE)
	return utils.Wrap(err)
}
