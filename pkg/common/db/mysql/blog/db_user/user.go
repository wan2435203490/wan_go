package db_user

import (
	"database/sql"
	"regexp"
	"strings"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/db/mysql/blog/db_im_chat_group_friend"
	"wan_go/pkg/common/db/mysql/blog/db_im_chat_group_user"
	"wan_go/pkg/common/db/mysql/blog/db_wei_yan"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
	"wan_go/sdk/api"
)

const REGEX = "\\d{11}"

func Register(userIn *blogVO.UserVO) *blogVO.UserVO {
	var userOut blogVO.UserVO
	if matched, err := regexp.MatchString(userIn.UserName, REGEX); !matched || err != nil {
		userOut.ErrMsg = "用户名不能为11位数字！"
		return &userOut
	}

	if strings.Contains(userIn.UserName, "@") {
		userOut.ErrMsg = "用户名不能包含@！"
		return &userOut
	}

	if utils.IsNotEmpty(userIn.PhoneNumber) && utils.IsNotEmpty(userIn.Email) {
		userOut.ErrMsg = "手机号与邮箱只能选择其中一个！"
		return &userOut
	}

	if utils.IsNotEmpty(userIn.PhoneNumber) {
		key := blog_const.FORGET_PASSWORD + userIn.PhoneNumber + "_1"
		get, b := cache.GetString(key)
		if !b || get != userIn.Code {
			userOut.ErrMsg = "验证码错误！"
			return &userOut
		}
		cache.Delete(key)
	} else if utils.IsNotEmpty(userIn.Email) {
		key := blog_const.FORGET_PASSWORD + userIn.Email + "_2"
		get, b := cache.GetString(key)
		if !b || get != userIn.Code {
			userOut.ErrMsg = "验证码错误！"
			return &userOut
		}
		cache.Delete(key)
	} else {
		userOut.ErrMsg = "请输入邮箱或手机号！"
		return &userOut
	}

	enc, _ := utils.AesEncrypt([]byte(userIn.Password), []byte(blog_const.CRYPOTJS_KEY))
	userIn.Password = string(enc)

	var count int64
	if err := db.Mysql().Model(&blog.User{}).Where("user_name=?", userIn.UserName).Count(&count).Error; err != nil {
		userOut.ErrMsg = err.Error()
		return &userOut
	}
	if count > 0 {
		userOut.ErrMsg = "用户名重复！"
		return &userOut
	}

	if utils.IsNotEmpty(userIn.PhoneNumber) {
		//下面就不去判错了
		db.Mysql().Model(&blog.User{}).Where("phone_number=?", userIn.PhoneNumber).Count(&count)
		if count > 0 {
			userOut.ErrMsg = "手机号重复！"
			return &userOut
		}
	} else if utils.IsNotEmpty(userIn.Email) {
		//下面就不去判错了
		db.Mysql().Model(&blog.User{}).Where("email=?", userIn.Email).Count(&count)
		if count > 0 {
			userOut.ErrMsg = "邮箱重复！"
			return &userOut
		}
	}

	user := blog.User{}
	user.UserName = userIn.UserName
	user.PhoneNumber = userIn.PhoneNumber
	user.Email = userIn.Email
	user.Password = utils.Md5(userIn.Password)

	if utils.IsEmpty(user.Avatar) {
		//todo
		//userOut.Avatar = randomavatar
	}

	if err := Insert(&user); err != nil {
		userOut.ErrMsg = err.Error()
		return &userOut
	}

	if err := db.Mysql().Find(&user).Error; err != nil {
		userOut.ErrMsg = err.Error()
		return &userOut
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
	weiYan.Content = "我也到此一游"
	weiYan.Type = blog_const.WEIYAN_TYPE_FRIEND
	weiYan.IsPublic = true
	_ = db_wei_yan.Insert(&weiYan)

	imChatGroupUser := blog.ImChatGroupUser{}
	imChatGroupUser.UserId = user.ID
	imChatGroupUser.GroupId = blog_const.DEFAULT_GROUP_ID
	imChatGroupUser.UserStatus = blog_const.GROUP_USER_STATUS_PASS
	_ = db_im_chat_group_user.Insert(&imChatGroupUser)

	imChatUser := blog.ImChatUserFriend{}
	imChatUser.UserId = user.ID
	imChatUser.FriendId = cache.GetAdminUser().ID
	imChatUser.Remark = "站长"
	imChatUser.FriendStatus = blog_const.FRIEND_STATUS_PASS
	_ = db_im_chat_group_friend.Insert(&imChatUser)

	imChatFriend := blog.ImChatUserFriend{}
	imChatFriend.UserId = cache.GetAdminUser().ID
	imChatFriend.FriendId = user.ID
	imChatFriend.FriendStatus = blog_const.FRIEND_STATUS_PASS
	_ = db_im_chat_group_friend.Insert(&imChatFriend)

	return &userOut
}

func Login(account string, password []byte, isAdmin bool) *blogVO.UserVO {

	var userVO blogVO.UserVO

	var err error
	if password, err = utils.AesDecrypt(password, []byte(blog_const.CRYPOTJS_KEY)); err != nil {
		return nil
	}

	var user blog.User
	if err = db.Mysql().Where("password = ? and (user_name = @account or email = @account or phone_number = @account)",
		password, sql.Named("account", account)).First(&user).Error; err != nil {
		userVO.ErrMsg = "账号/密码错误，请重新输入！"
		return &userVO
	}

	if !user.UserStatus {
		userVO.ErrMsg = "账号被冻结！"
		return &userVO
	}

	if isAdmin {
		adminLogin(&user, &userVO)
	} else {
		userLogin(&user, &userVO)
	}

	userVO.Copy(&user)
	userVO.Password = ""

	return &userVO
}

func adminLogin(user *blog.User, userVO *blogVO.UserVO) {
	var adminToken string

	if user.UserType != blog_const.USER_TYPE_ADMIN.Code && user.UserType != blog_const.USER_TYPE_DEV.Code {
		userVO.ErrMsg = "请输入管理员账号！"
	}

	key := blog_const.ADMIN_TOKEN + utils.Int32ToString(user.ID)
	if get, b := cache.Get(key); b {
		adminToken = get.(string)
	}

	if utils.IsEmpty(adminToken) {
		adminToken = blog_const.ADMIN_ACCESS_TOKEN + utils.UUID()
		cache.Set(adminToken, &user)
		cache.Set(key, adminToken)
	}

	if user.UserType == blog_const.USER_TYPE_ADMIN.Code {
		userVO.IsBoss = true
	}
	userVO.AccessToken = adminToken
}

func userLogin(user *blog.User, userVO *blogVO.UserVO) {
	var userToken string

	key := blog_const.USER_TOKEN + utils.Int32ToString(user.ID)
	if get, b := cache.Get(key); b {
		userToken = get.(string)
	}

	if utils.IsEmpty(userToken) {
		userToken = blog_const.USER_ACCESS_TOKEN + utils.UUID()
		cache.Set(userToken, &user)
		cache.Set(key, userToken)
	}
	userVO.AccessToken = userToken
}

func LoginByToken(token string) *blogVO.UserVO {

	var userVO blogVO.UserVO
	token, err := utils.AesDecryptByString(token, blog_const.CRYPOTJS_KEY)
	if err != nil || utils.IsEmpty(token) {
		userVO.ErrMsg = "未登录，请登陆后再进行操作！"
		return &userVO
	}
	if get, b := cache.Get(token); !b {
		userVO.ErrMsg = "登录已过期，请重新登录！"
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

func UpdateUserInfo(userIn *blogVO.UserVO, userToken string) *blogVO.UserVO {
	var userOut blogVO.UserVO
	if matched, err := regexp.MatchString(userIn.UserName, REGEX); !matched || err != nil {
		userOut.ErrMsg = "用户名不能为11位数字！"
		return &userOut
	}

	if strings.Contains(userIn.UserName, "@") {
		userOut.ErrMsg = "用户名不能包含@！"
		return &userOut
	}

	var count int64
	if err := db.Mysql().Model(&blog.User{}).
		Where("user_name=? and id <> ?", userIn.UserName, userIn.ID).
		Count(&count).Error; err != nil {
		userOut.ErrMsg = err.Error()
		return &userOut
	}
	if count > 0 {
		userOut.ErrMsg = "用户名重复！"
		return &userOut
	}

	user := blog.User{}
	user.ID = userIn.ID
	user.UserName = userIn.UserName
	user.Avatar = userIn.Avatar
	user.Gender = userIn.Gender
	user.Introduction = userIn.Introduction
	if err := Update(&user); err != nil {
		userOut.ErrMsg = err.Error()
		return &userOut
	}

	if err := db.Mysql().Find(&user).Error; err != nil {
		userOut.ErrMsg = err.Error()
		return &userOut
	}

	key := blog_const.USER_TOKEN + utils.Int32ToString(user.ID)
	cache.Set(userToken, &user)
	cache.Set(key, userToken)

	userOut.Copy(&user)
	userOut.Password = ""

	return &userOut
}

func UpdateSecretInfo(place, flag, captcha, password string, user *blog.User) *blogVO.UserVO {
	password, _ = utils.AesDecryptByString(password, blog_const.CRYPOTJS_KEY)

	userVO := blogVO.UserVO{}
	if flag == "1" || flag == "2" {
		if utils.Md5(password) != user.Password {
			userVO.ErrMsg = "密码错误！"
			return &userVO
		}

		if utils.IsEmpty(captcha) {
			userVO.ErrMsg = "请输入验证码！"
			return &userVO
		}
	}

	updateUser := blog.User{}
	updateUser.ID = user.ID

	var count int64
	//todo 统一管理key
	key := blog_const.USER_CODE + utils.Int32ToString(user.ID) + "_" + place + "_" + flag

	switch flag {
	case "1":
		if err := db.Mysql().Model(&blog.User{}).Where("phone_number=?", place).Count(&count).Error; err != nil {
			//todo 包装errmsg
			userVO.ErrMsg = err.Error()
			return &userVO
		}
		if count > 0 {
			userVO.ErrMsg = "手机号重复！"
			return &userVO
		}

		fun := func() { updateUser.PhoneNumber = place }
		if !validateCaptcha(key, captcha, &userVO, fun) {
			return &userVO
		}

	case "2":
		if err := db.Mysql().Model(&blog.User{}).Where("email=?", place).Count(&count).Error; err != nil {
			userVO.ErrMsg = err.Error()
			return &userVO
		}
		if count > 0 {
			userVO.ErrMsg = "邮箱重复！"
			return &userVO
		}
		fun := func() { updateUser.Email = place }
		if !validateCaptcha(key, captcha, &userVO, fun) {
			return &userVO
		}
	case "3":
		if utils.Md5(place) == user.Password {
			updateUser.Password = utils.Md5(password)
		} else {
			userVO.ErrMsg = "密码错误！"
			return &userVO
		}
	default:
		break
	}

	if err := Update(&updateUser); err != nil {
		userVO.ErrMsg = err.Error()
		return &userVO
	}
	if err := db.Mysql().Find(&updateUser).Error; err != nil {
		userVO.ErrMsg = err.Error()
		return &userVO
	}

	key = blog_const.USER_TOKEN + utils.Int32ToString(user.ID)
	cache.Set(utils.Token(), &updateUser)
	cache.Set(key, utils.Token())

	userVO.Copy(&updateUser)
	userVO.Password = ""
	return &userVO
}

func validateCaptcha(key, captcha string, userVO *blogVO.UserVO, fun func()) bool {
	captchaCache, ok := cache.GetString(key)
	if ok && captchaCache == captcha {
		cache.Delete(key)
		fun()
		return true
	} else {
		userVO.ErrMsg = "验证码错误！"
		return false
	}
}

func WrapError(msg string) *api.Err {
	return &api.Err{ErrMsg: msg}
}

func UpdateForForgetPassword(place, flag, captcha, password string) *api.Err {

	password, _ = utils.AesDecryptByString(password, blog_const.CRYPOTJS_KEY)

	key := blog_const.FORGET_PASSWORD + place + "_" + flag
	codeCache, b := cache.GetString(key)
	if !b || codeCache != captcha {
		return WrapError("验证码错误！")
	}

	cache.Delete(key)

	newPassword := utils.Md5(password)

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

//db.DB.MysqlDB.Where("name LIKE ?", "group%")

func ListWith(query interface{}, args ...interface{}) ([]*blog.User, error) {

	var users []*blog.User

	tx := db.Mysql()

	if query == nil || args == nil {
		tx = tx.Select(query, args)
	}

	if err := tx.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func List() ([]*blog.User, error) {
	return ListWith(nil, nil)
}

const LimitUserList = 10

func ListByUserName(userName string) []*blogVO.UserVO {

	var users []*blog.User
	if err := db.Mysql().Select("id, user_name, avatar, gender, introduction").
		Where("user_name like ?", userName+"%").
		Limit(LimitUserList).Last(&users).Error; err != nil {
		return nil
	}

	var result = make([]*blogVO.UserVO, LimitUserList)
	for _, user := range users {
		vo := blogVO.UserVO{}
		vo.Copy(user)
		result = append(result, &vo)
	}

	return result
}

func ListUser(vo *blogVO.BaseRequestVO[*blog.User]) {

	var users []*blog.User
	tx := db.Page(&vo.Pagination).Where("user_status=?", vo.UserStatus)
	if vo.UserType > 0 {
		tx.Where("user_type=?", vo.UserType)
	}
	if utils.IsNotEmpty(vo.SearchKey) {
		tx.Where("(user_name=@searchKey or phone_number=@searchKey)", sql.Named("searchKey", vo.SearchKey))
	}
	if err := tx.Omit("password, open_id").Order("CreatedAt DESC").Find(&users).Error; err != nil {
		return
	}

	vo.Total = len(users)
	vo.Records = users
}

func Update(user *blog.User) error {
	return db.Mysql().Updates(&user).Error
}

func Insert(user *blog.User) error {
	return db.Mysql().Create(&user).Error
}

func Delete(user *blog.User) error {
	return db.Mysql().Delete(&user).Error
}
