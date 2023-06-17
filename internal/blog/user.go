package blog

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog/db_user"
	"wan_go/pkg/common/log"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
)

func Register(c *gin.Context) {
	var vo blogVO.UserVO
	if a.BindFailed(&vo) {
		return
	}

	userVO := db_user.Register(&vo)

	a.OK(userVO)
}

func Login(c *gin.Context) {

	account := a.Param("account")
	password := a.Param("password")
	//isAdmin := a.Param("isAdmin")

	userVO := db_user.Login(account, []byte(password), false)

	a.OK(userVO)
}

func LoginByToken(c *gin.Context) {

	userToken := a.Param("userToken")

	userVO := db_user.LoginByToken(userToken)

	a.OK(userVO)
}

func Logout(c *gin.Context) {
	token := a.GetToken()
	userId := a.GetCurrentUserId()
	db_user.Exit(token, userId)
	a.OK()
}

func UpdateUserInfo(c *gin.Context) {
	var vo blogVO.UserVO
	if a.BindFailed(&vo) {
		return
	}

	cache.Delete(a.KeyUserId(blog_const.USER_CACHE))

	vo.ID = a.GetCurrentUserId()
	userToken := a.GetToken()
	db_user.UpdateUserInfo(&vo, userToken)
	a.OK(&vo)
}

/**
 * 获取验证码
 * <p>
 * 1 手机号
 * 2 邮箱
 */
func GetCaptcha(c *gin.Context) {

	flag := a.Param("flag")
	captcha := utils.CreateCaptcha(6)
	user := a.GetCurrentUser()

	switch flag {
	case "1":
		if a.EmptyFailed("请先绑定手机号！", user.PhoneNumber) {
			return
		}
		log.Info("GetCaptcha", user.ID, "---"+user.UserName+"---"+"手机验证码:"+captcha)
	case "2":
		if a.EmptyFailed("请先绑定邮箱！", user.Email) {
			return
		}
		log.Info("GetCaptcha", user.ID, "---"+user.UserName+"---"+"邮箱验证码:"+captcha)
		sendMail(captcha, user.Email)
	default:
		break
	}

	cache.SetExpire(a.KeyUserId(blog_const.USER_CODE)+"_"+flag, captcha, time.Minute*5)

	a.OK()
}

// GetCaptchaForBind
// 绑定手机号或者邮箱
// flag: 1 手机号 2 邮箱
func GetCaptchaForBind(c *gin.Context) {
	place := a.Param("place")
	flag := a.Param("flag")

	captcha := utils.CreateCaptcha(6)
	user := a.GetCurrentUser()

	switch flag {
	case "1":
		log.Info("GetCodeForBind", place+"---"+"手机验证码:"+captcha)
	case "2":
		log.Info("GetCodeForBind", place+"---"+"邮箱验证码:"+captcha)
		sendMail(captcha, user.Email)
	default:
		break
	}

	cache.SetExpire(a.KeyUserId(blog_const.USER_CODE)+"_"+place+"_"+flag, captcha, time.Minute*5)

	a.OK()
}

// UpdateSecretInfo
// 更新邮箱、手机号、密码
// 1 手机号
// 2 邮箱
// 3 密码：place=老密码&password=新密码
func UpdateSecretInfo(c *gin.Context) {
	place := a.Param("place")
	flag := a.Param("flag")
	captcha := a.Param("code")
	password := a.Param("password")

	cache.Delete(a.KeyUserId(blog_const.USER_CACHE))

	user := a.GetCurrentUser()

	a.OK(db_user.UpdateSecretInfo(place, flag, captcha, password, user))
}

// GetCaptchaForForgetPassword
// 忘记密码 获取验证码
// 1 手机号
// 2 邮箱
func GetCaptchaForForgetPassword(c *gin.Context) {
	place := a.Param("place")
	flag := a.Param("flag")

	captcha := utils.CreateCaptcha(6)

	switch flag {
	case "1":
		log.Info("GetCaptchaForForgetPassword", "手机验证码:"+captcha)
	case "2":
		log.Info("GetCaptchaForForgetPassword", "邮箱验证码:"+captcha)
		sendMail(captcha, place)
	default:
		break
	}

	cache.SetExpire(blog_const.FORGET_PASSWORD+place+"_"+flag, captcha, time.Minute*5)

	a.OK()
}

// UpdateForForgetPassword
// 忘记密码 更新密码
// 1 手机号
// 2 邮箱
func UpdateForForgetPassword(c *gin.Context) {
	place := a.Param("place")
	flag := a.Param("flag")
	captcha := a.Param("code")
	password := a.Param("password")

	apiErr := db_user.UpdateForForgetPassword(place, flag, captcha, password)

	if a.IsFailed(apiErr != nil, apiErr.ErrMsg) {
		return
	}

	a.OK()
}

func GetUserByUsername(c *gin.Context) {
	userName := a.Param("username")
	a.OK(db_user.ListByUserName(userName))
}

func sendMail(captcha, email string) {
	mails := make([]string, 0)
	mails = append(mails, email)
	webName := cache.GetWebName()
	codeMail := getCodeMail(captcha, webName)

	utils.SendMail(mails, "您有一封来自"+webName+"的回执！", codeMail)
}

func getCodeMail(captcha, webName string) string {

	admin := cache.GetAdminUser()
	var userName string
	if admin == nil {
		userName = config.Config.Blog.Name
	} else {
		userName = admin.UserName
	}
	codeMail := fmt.Sprintf(utils.MailText,
		webName,
		fmt.Sprintf(utils.ImMail, userName),
		userName,
		fmt.Sprintf(config.Config.User.CaptchaFormat, captcha),
		"",
		webName)

	return codeMail
}
