package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog/db_user"
	"wan_go/pkg/common/log"
	"wan_go/pkg/common/mail"
	"wan_go/pkg/utils"
	blogVO "wan_go/pkg/vo/blog"
)

type UserApi struct {
	api.Api
}

func (a UserApi) Register(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.UserVO
	if a.BindFailed(&vo) {
		return
	}

	if a.IsFailed(utils.IsEmpty(vo.UserName), "用户名不能为空") {
		return
	}

	if a.IsFailed(utils.IsEmpty(vo.Password), "密码不能为空") {
		return
	}

	userVO := db_user.Register(&vo)

	a.OK(userVO)
}

func (a UserApi) Login(c *gin.Context) {
	a.MakeContext(c)

	var loginVO blogVO.LoginVO
	if a.BindFailed(&loginVO) {
		return
	}
	if a.IsFailed(utils.IsEmpty(loginVO.Password), "密码不能为空") {
		return
	}

	userVO := db_user.Login(loginVO.Account, loginVO.Password, loginVO.IsAdmin)

	a.OK(userVO)
}

func (a UserApi) LoginByToken(c *gin.Context) {
	a.MakeContext(c)

	userToken := a.Param("userToken")

	userVO := db_user.LoginByToken(userToken)

	a.OK(userVO)
}

func (a UserApi) Logout(c *gin.Context) {
	a.MakeContext(c)
	token := a.GetToken()
	userId := a.GetCurrentUserId()
	db_user.Exit(token, userId)
	a.OK()
}

func (a UserApi) UpdateUserInfo(c *gin.Context) {
	a.MakeContext(c)
	var vo blogVO.UserVO
	if a.BindFailed(&vo) {
		return
	}

	cache.Delete(a.KeyUserId(blog_const.USER_CACHE))

	vo.ID = a.GetCurrentUserId()
	userToken := a.GetToken()
	ret := db_user.UpdateUserInfo(&vo, userToken)
	a.OK(ret)
}

/**
 * 获取验证码
 * <p>
 * 1 手机号
 * 2 邮箱
 */
func (a UserApi) GetCaptcha(c *gin.Context) {
	a.MakeContext(c)

	flag := a.QueryInt("flag")

	captcha := utils.CreateCaptcha(6)
	user := a.GetCurrentUser()

	switch flag {
	case 1:
		if a.EmptyFailed("请先绑定手机号！", user.PhoneNumber) {
			return
		}
		log.Info("GetCaptcha", user.ID, "---"+user.UserName+"---"+"手机验证码:"+captcha)
	case 2:
		if a.EmptyFailed("请先绑定邮箱！", user.Email) {
			return
		}
		log.Info("GetCaptcha", user.ID, "---"+user.UserName+"---"+"邮箱验证码:"+captcha)
		sendMail(captcha, user.Email)
	default:
		break
	}

	cache.SetExpire(a.KeyUserId(blog_const.USER_CODE)+"_"+utils.IntToString(flag), captcha, time.Minute*5)

	a.OK()
}

// GetCaptchaForBind
// 绑定手机号或者邮箱
// flag: 1 手机号 2 邮箱
func (a UserApi) GetCaptchaForBind(c *gin.Context) {
	a.MakeContext(c)
	//place := a.Param("place")
	//flag := a.Param("flag")
	var param blogVO.Param
	if a.BindFailed(&param) {
		return
	}

	captcha := utils.CreateCaptcha(6)
	user := a.GetCurrentUser()

	switch param.Flag {
	case 1:
		log.Info("GetCodeForBind", param.Place+"---"+"手机验证码:"+captcha)
	case 2:
		log.Info("GetCodeForBind", param.Place+"---"+"邮箱验证码:"+captcha)
		sendMail(captcha, user.Email)
	default:
		break
	}

	cache.SetExpire(a.KeyUserId(blog_const.USER_CODE)+"_"+param.Place+"_"+param.FlagString(), captcha, time.Minute*5)

	a.OK()
}

// UpdateSecretInfo
// 更新邮箱、手机号、密码
// 1 手机号
// 2 邮箱
// 3 密码：place=老密码&password=新密码
func (a UserApi) UpdateSecretInfo(c *gin.Context) {
	a.MakeContext(c)
	//place := a.Param("place")
	//flag := a.Param("flag")
	//captcha := a.Param("code")
	//password := a.Param("password")
	var param blogVO.Param
	if a.BindFailed(&param) {
		return
	}

	cache.Delete(a.KeyUserId(blog_const.USER_CACHE))

	user := a.GetCurrentUser()

	a.OK(db_user.UpdateSecretInfo(param.Place, param.FlagString(), param.Code, param.Password, user))
}

// GetCaptchaForForgetPassword
// 忘记密码 获取验证码
// 1 手机号
// 2 邮箱
func (a UserApi) GetCaptchaForForgetPassword(c *gin.Context) {
	a.MakeContext(c)
	place := a.Query("place")
	flag := a.QueryInt("flag")
	param := blogVO.Param{Place: place, Flag: flag}

	captcha := utils.CreateCaptcha(6)

	switch param.Flag {
	case 1:
		log.Info("GetCaptchaForForgetPassword", "手机验证码:"+captcha)
	case 2:
		log.Info("GetCaptchaForForgetPassword", "邮箱验证码:"+captcha)
		sendMail(captcha, param.Place)
	default:
		break
	}

	cache.SetExpire(blog_const.FORGET_PASSWORD+param.Place+"_"+param.FlagString(), captcha, time.Minute*5)

	a.OK()
}

// UpdateForForgetPassword
// 忘记密码 更新密码
// 1 手机号
// 2 邮箱
func (a UserApi) UpdateForForgetPassword(c *gin.Context) {
	a.MakeContext(c)
	place := a.Param("place")
	flag := a.Param("flag")
	captcha := a.Param("code")
	password := a.Param("password")

	apiErr := db_user.UpdateForForgetPassword(place, flag, captcha, password)

	if a.IsFailed(apiErr != nil, apiErr.Msg) {
		return
	}

	a.OK()
}

func (a UserApi) GetUserByUsername(c *gin.Context) {
	a.MakeContext(c)
	userName := a.Param("username")
	a.OK(db_user.ListByUserName(userName))
}

func sendMail(captcha, email string) {
	mails := make([]string, 0)
	mails = append(mails, email)
	webName := cache.GetWebName()
	codeMail := getCodeMail(captcha, webName)

	mail.SendMail(mails, "您有一封来自"+webName+"的回执！", codeMail)
}

func getCodeMail(captcha, webName string) string {

	admin := cache.GetAdminUser()
	var userName string
	if admin == nil {
		userName = config.Config.Blog.Name
	} else {
		userName = admin.UserName
	}
	codeMail := fmt.Sprintf(mail.MailText,
		webName,
		fmt.Sprintf(mail.ImMail, userName),
		userName,
		fmt.Sprintf(config.Config.User.CaptchaFormat, captcha),
		"",
		webName)

	return codeMail
}
