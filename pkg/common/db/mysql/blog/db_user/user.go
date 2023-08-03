package db_user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
	"wan_go/internal/blog/service"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
	r "wan_go/pkg/common/response"
	"wan_go/pkg/utils"
)

// 匹配规则// ^1第一位为一// [345789]{1} 后接一位345789 的数字// \\d \d的转义 表示数字 {9} 接9位
const RegRuler = "^1[345789]{1}\\d{9}$"

// todo 新增正则util
func isMobile(phoneNumber string) bool {
	//正则调用规则
	reg := regexp.MustCompile(RegRuler)
	// 返回 MatchString 是否匹配
	return reg.MatchString(phoneNumber)
}

func Register(c *gin.Context, userIn *vo.UserVO) (*vo.UserVO, error) {
	var userOut vo.UserVO
	if isMobile(userIn.UserName) {
		return nil, errors.New("用户名不能为电话号码！")
	}

	if strings.Contains(userIn.UserName, "@") {
		return nil, errors.New("用户名不能包含@！")
	}

	if utils.IsNotEmpty(userIn.PhoneNumber) && utils.IsNotEmpty(userIn.Email) {
		return nil, errors.New("手机号与邮箱只能选择其中一个！")
	}

	if utils.IsNotEmpty(userIn.PhoneNumber) {
		key := blog_const.FORGET_PASSWORD + userIn.PhoneNumber + "_1"
		get, b := cache.GetString(key)
		if !b || get != userIn.Captcha {
			return nil, errors.New("验证码错误！")
		}
		cache.Delete(key)
	} else if utils.IsNotEmpty(userIn.Email) {
		key := blog_const.FORGET_PASSWORD + userIn.Email + "_2"
		get, b := cache.GetString(key)
		if !b || get != userIn.Captcha {
			return nil, errors.New("验证码错误！")
		}
		cache.Delete(key)
	} else {
		return nil, errors.New("请输入邮箱或手机号！")
	}

	jsEncodePwd := userIn.Password
	dec := utils.AesDecryptCrypotJsKey(userIn.Password)
	userIn.Password = dec

	var count int64
	if err := db.Mysql().Model(&blog.User{}).Where("user_name=?", userIn.UserName).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名重复！")
	}

	if utils.IsNotEmpty(userIn.PhoneNumber) {
		//下面就不去判错了
		db.Mysql().Model(&blog.User{}).Where("phone_number=?", userIn.PhoneNumber).Count(&count)
		if count > 0 {
			return nil, errors.New("手机号重复！")
		}
	} else if utils.IsNotEmpty(userIn.Email) {
		//下面就不去判错了
		db.Mysql().Model(&blog.User{}).Where("email=?", userIn.Email).Count(&count)
		if count > 0 {
			return nil, errors.New("邮箱重复！")
		}
	}

	user := blog.User{}
	user.UserName = userIn.UserName
	user.PhoneNumber = userIn.PhoneNumber
	user.Email = userIn.Email
	//user.Password = utils.Md5(userIn.Password)
	user.Password = userIn.Password
	user.CrypotJsText = jsEncodePwd

	if utils.IsEmpty(user.Avatar) {
		user.Avatar = utils.RandomAvatar(0)
	}

	if err := Insert(&user); err != nil {
		return nil, err
	}

	if err := db.Mysql().Find(&user).Error; err != nil {
		return nil, err
	}

	userToken := blog_const.USER_ACCESS_TOKEN + utils.UUID()
	key := blog_const.USER_TOKEN + utils.Int32ToString(user.ID)
	cache.Set(userToken, &user)
	cache.Set(key, userToken)

	userOut.Copy(&user)
	userOut.Password = ""
	userOut.AccessToken = userToken

	weiYan := blog.WeiYan{}
	weiYan.UserId = user.ID
	weiYan.Content = "到此一游~"
	weiYan.Type = blog_const.WEIYAN_TYPE_FRIEND
	weiYan.IsPublic = true
	wys := service.NewWeiYan(c)
	_ = wys.Insert(&weiYan)

	/** im
	//
	//imChatGroupUser := blog.ImChatGroupUser{}
	//imChatGroupUser.UserId = user.ID
	//imChatGroupUser.GroupId = blog_const.DEFAULT_GROUP_ID
	//imChatGroupUser.UserStatus = blog_const.GROUP_USER_STATUS_PASS
	//_ = db_im_chat_group_user.InsertReq(&imChatGroupUser)
	//
	//imChatUser := blog.ImChatUserFriend{}
	//imChatUser.UserId = user.ID
	//imChatUser.FriendId = int32(cache.GetAdminUserId())
	//imChatUser.Remark = "站长"
	//imChatUser.FriendStatus = blog_const.FRIEND_STATUS_PASS
	//_ = db_im_chat_group_friend.InsertReq(&imChatUser)
	//
	//imChatFriend := blog.ImChatUserFriend{}
	//imChatFriend.UserId = int32(cache.GetAdminUserId())
	//imChatFriend.FriendId = user.ID
	//imChatFriend.FriendStatus = blog_const.FRIEND_STATUS_PASS
	//_ = db_im_chat_group_friend.InsertReq(&imChatFriend)
	*/

	return &userOut, nil
}

func LoginByToken(token string) *vo.UserVO {

	var userVO vo.UserVO
	token = utils.AesDecryptCrypotJsKey(token)
	if utils.IsEmpty(token) {
		userVO.Msg = "未登录，请登录后再进行操作！"
		return &userVO
	}
	if get, b := cache.Get(token); !b {
		userVO.Msg = "登录已过期，请重新登录！"
		return &userVO
	} else {
		userVO.Copy(get.(*blog.User))
		userVO.Password = ""
		userVO.AccessToken = token
	}
	return &userVO
}

func Exit(token string, userId int32) {
	if strings.Contains(token, blog_const.USER_ACCESS_TOKEN) {
		cache.Delete(blog_const.USER_TOKEN + utils.Int32ToString(userId))
		//todo im
		//websocket 移除连接
	} else if strings.Contains(token, blog_const.ADMIN_ACCESS_TOKEN) {
		cache.Delete(blog_const.ADMIN_TOKEN + utils.Int32ToString(userId))
	}
	cache.Delete(token)
}

func UpdateUserInfo(userIn *vo.UserVO) *vo.UserVO {
	var userOut vo.UserVO
	if isMobile(userIn.UserName) {
		userOut.Msg = "用户名不能为电话号码！"
		return &userOut
	}

	if strings.Contains(userIn.UserName, "@") {
		userOut.Msg = "用户名不能包含@！"
		return &userOut
	}

	var count int64
	if err := db.Mysql().Model(&blog.User{}).
		Where("user_name=? and id <> ?", userIn.UserName, userIn.ID).
		Count(&count).Error; err != nil {
		userOut.Msg = err.Error()
		return &userOut
	}
	if count > 0 {
		userOut.Msg = "用户名重复！"
		return &userOut
	}

	user := blog.User{}
	user.ID = userIn.ID
	user.UserName = userIn.UserName
	user.Avatar = userIn.Avatar
	user.Gender = userIn.Gender
	user.Introduction = userIn.Introduction
	if err := Update(&user); err != nil {
		userOut.Msg = err.Error()
		return &userOut
	}

	if err := db.Mysql().Find(&user).Error; err != nil {
		userOut.Msg = err.Error()
		return &userOut
	}

	userOut.Copy(&user)
	userOut.Password = ""

	return &userOut
}

func UpdateSecretInfo(place, flag, captcha, password string, userId int32) *vo.UserVO {
	password = utils.AesDecryptCrypotJsKey(password)

	userVO := vo.UserVO{}
	if flag == "1" || flag == "2" {
		//token校验了
		//if utils.Md5(password) != user.Password {
		//	userVO.Message = "密码错误！"
		//	return &userVO
		//}

		if utils.IsEmpty(captcha) {
			userVO.Msg = "请输入验证码！"
			return &userVO
		}
	}

	updateUser := blog.User{}
	updateUser.ID = userId

	var count int64
	//todo 统一管理key
	key := blog_const.USER_CODE + utils.Int32ToString(userId) + "_" + place + "_" + flag

	switch flag {
	case "1":
		if err := db.Mysql().Model(&blog.User{}).Where("phone_number=?", place).Count(&count).Error; err != nil {
			//todo 包装errmsg
			userVO.Msg = err.Error()
			return &userVO
		}
		if count > 0 {
			userVO.Msg = "手机号重复！"
			return &userVO
		}

		fun := func() { updateUser.PhoneNumber = place }
		if !validateCaptcha(key, captcha, &userVO, fun) {
			return &userVO
		}

	case "2":
		if err := db.Mysql().Model(&blog.User{}).Where("email=?", place).Count(&count).Error; err != nil {
			userVO.Msg = err.Error()
			return &userVO
		}
		if count > 0 {
			userVO.Msg = "邮箱重复！"
			return &userVO
		}
		fun := func() { updateUser.Email = place }
		if !validateCaptcha(key, captcha, &userVO, fun) {
			return &userVO
		}
	case "3":
		//if utils.Md5(place) == user.Password {
		//	updateUser.Password = utils.Md5(password)
		//} else {
		//	userVO.Message = "密码错误！"
		//	return &userVO
		//}
	default:
		break
	}

	if err := Update(&updateUser); err != nil {
		userVO.Msg = err.Error()
		return &userVO
	}
	if err := db.Mysql().Find(&updateUser).Error; err != nil {
		userVO.Msg = err.Error()
		return &userVO
	}

	key = blog_const.USER_TOKEN + utils.Int32ToString(userId)
	cache.Set(utils.Token(), &updateUser)
	cache.Set(key, utils.Token())

	userVO.Copy(&updateUser)
	userVO.Password = ""
	return &userVO
}

func validateCaptcha(key, captcha string, userVO *vo.UserVO, fun func()) bool {
	captchaCache, ok := cache.GetString(key)
	if ok && captchaCache == captcha {
		cache.Delete(key)
		fun()
		return true
	} else {
		userVO.Msg = "验证码错误！"
		return false
	}
}

func WrapError(msg string) *r.CodeMsg {
	return &r.CodeMsg{Msg: msg}
}

func UpdateForForgetPassword(place, flag, captcha, password string) *r.CodeMsg {

	password = utils.AesDecryptCrypotJsKey(password)

	key := blog_const.FORGET_PASSWORD + place + "_" + flag
	codeCache, b := cache.GetString(key)
	if !b || codeCache != captcha {
		return WrapError("验证码错误！")
	}

	cache.Delete(key)

	//newPassword := utils.Md5(password)
	newPassword := password

	switch flag {
	case "1":
		user := blog.User{}
		if err := db.Mysql().Where("phone_number=?", place).Find(&user).Error; err != nil {
			return WrapError(err.Error())
		}
		if user.ID == 0 {
			return WrapError("该手机号未绑定账号！")
		}
		if !user.UserStatus {
			return WrapError("账号被冻结！")
		}
		db.Mysql().Model(&blog.User{}).Where("phone_number=?", place).Update("password", newPassword)
		cache.Delete(blog_const.USER_CACHE + utils.Int32ToString(user.ID))

	case "2":
		user := blog.User{}
		if err := db.Mysql().Where("email=?", place).Find(&user).Error; err != nil {
			return WrapError(err.Error())
		}
		if user.ID == 0 {
			return WrapError("该邮箱未绑定账号！")
		}
		if !user.UserStatus {
			return WrapError("账号被冻结！")
		}
		db.Mysql().Model(&blog.User{}).Where("email=?", place).Update("password", newPassword)
		cache.Delete(blog_const.USER_CACHE + utils.Int32ToString(user.ID))

	default:
		break
	}

	return nil
}

func Get(user *blog.User) error {
	return db.Mysql().Find(&user).Error
}

const LimitUserList = 10

func ListByUserName(userName string) []*vo.UserVO {

	var users []*blog.User
	if err := db.Mysql().Select("id, user_name, avatar, gender, introduction").
		Where("user_name like ?", userName+"%").
		Limit(LimitUserList).Last(&users).Error; err != nil {
		return nil
	}

	var result = make([]*vo.UserVO, LimitUserList)
	for _, user := range users {
		vo := vo.UserVO{}
		vo.Copy(user)
		result = append(result, &vo)
	}

	return result
}

func Update(user *blog.User) error {
	return db.Mysql().Updates(&user).Error
}

func Insert(user *blog.User) error {
	return db.Mysql().Create(&user).Error
}
