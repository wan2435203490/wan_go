package service

import (
	"github.com/gin-gonic/gin"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/sdk/pkg"
	"wan_go/sdk/service"
)

type User struct {
	service.Service
}

func NewUser(c *gin.Context) *User {
	us := User{}
	us.Orm = pkg.Orm(c)
	us.Log = api.GetRequestLogger(c)
	return &us
}

//func (s *User) Login(account string, password string, isAdmin bool) *vo.UserVO {
//
//	var userVO vo.UserVO
//
//	var err error
//	password = utils.AesDecryptCrypotJsKey(password)
//
//	var user blog.User
//	if err = db.Mysql().Debug().Where("password = @password and email = @account)",
//		sql.Named("password", password), sql.Named("account", account)).First(&user).Error; err != nil {
//		userVO.Msg = "账号/密码错误，请重新输入！"
//		return &userVO
//	}
//
//	if !user.UserStatus {
//		userVO.Msg = "账号被冻结！"
//		return &userVO
//	}
//
//	if isAdmin {
//		adminLogin(&user, &userVO)
//	} else {
//		userLogin(&user, &userVO)
//	}
//
//	userVO.Copy(&user)
//	userVO.Password = ""
//
//	return &userVO
//}
//
//func adminLogin(user *blog.User, userVO *vo.UserVO) {
//	var token string
//
//	if user.UserType != blog_const.USER_TYPE_ADMIN.Code && user.UserType != blog_const.USER_TYPE_DEV.Code {
//		userVO.Msg = "请输入管理员账号！"
//	}
//
//	key := blog_const.ADMIN_TOKEN + utils.Int32ToString(user.ID)
//	if get, b := cache.Get(key); b {
//		token = get.(string)
//	}
//
//	if utils.IsEmpty(token) {
//		token = blog_const.ADMIN_ACCESS_TOKEN + utils.UUID()
//		cache.Set(token, user)
//		cache.Set(key, token)
//	}
//
//	if user.UserType == blog_const.USER_TYPE_ADMIN.Code {
//		userVO.IsBoss = true
//	}
//	userVO.AccessToken = token
//}
//
//func userLogin(user *blog.User, userVO *vo.UserVO) {
//	var token string
//
//	key := blog_const.USER_TOKEN + utils.Int32ToString(user.ID)
//	if get, b := cache.Get(key); b {
//		token = get.(string)
//	}
//
//	if utils.IsEmpty(token) {
//		token = blog_const.USER_ACCESS_TOKEN + utils.UUID()
//		cache.Set(token, user)
//		cache.Set(key, token)
//	}
//	userVO.AccessToken = token
//}

func (s *User) GetUser(user *blog.User) error {

	if err := s.Orm.Debug().Omit("password").Find(user).Error; err != nil {
		s.Log.Errorf("GetUser error: %s", err)
	}

	return nil
}

func (s *User) GetRole(role *blog.Role) error {

	if err := s.Orm.Debug().Find(role).Error; err != nil {
		s.Log.Errorf("GetRole error: %s", err)
	}

	return nil
}

func (s *User) BatchGetUser(users *[]blog.User) error {

	var ids []int32
	for _, user := range *users {
		ids = append(ids, user.ID)
	}

	if err := s.Orm.Debug().Model(&blog.User{}).Where("id in ?", ids).
		Find(users).Error; err != nil {
		s.Log.Errorf("GetUser error: %s", err)
		return err
	}

	return nil
}
