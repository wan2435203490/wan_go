package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"wan_go/internal/blog/service"
	"wan_go/internal/blog/vo"
	"wan_go/pkg/common/api"
	"wan_go/pkg/common/cache"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/db/mysql/blog/db_user"
	"wan_go/pkg/common/log"
	"wan_go/pkg/common/mail"
	"wan_go/pkg/utils"
	"wan_go/sdk/pkg/captcha"
	"wan_go/sdk/pkg/jwtauth/user"
)

type UserApi struct {
	api.Api
}

func (a UserApi) Register(c *gin.Context) {
	a.MakeContext(c)
	var user vo.UserVO
	if a.Bind(&user) {
		return
	}

	if a.IsFailed(utils.IsEmpty(user.UserName), "用户名不能为空") {
		return
	}

	if a.IsFailed(utils.IsEmpty(user.Password), "密码不能为空") {
		return
	}

	userVO, err := db_user.Register(c, &user)
	if a.IsError(err) {
		return
	}

	a.OK(userVO)
}

func (a UserApi) GetCaptchaImg(c *gin.Context) {
	a.MakeContext(c)

	id, b64s, err := captcha.DriverDigitFunc()
	if a.IsError(err) {
		return
	}
	a.Custom(gin.H{
		"code": 200,
		"data": b64s,
		"id":   id,
		"msg":  "success",
	})
}

func (a UserApi) Login(c *gin.Context) {
	//everything is in jwt middleware

	//a.MakeContext(c)
	//
	//var loginVO model.LoginVO
	//if a.Bind(&loginVO) {
	//	return
	//}
	//if a.IsFailed(utils.IsEmpty(loginVO.Password), "密码不能为空") {
	//	return
	//}
	//
	//userVO := db_user.Login(loginVO.Account, loginVO.Password, loginVO.IsAdmin)
	//
	//a.OK(userVO)
}

func (a UserApi) LoginByToken(c *gin.Context) {
	a.MakeContext(c)

	userToken := a.Param("userToken")

	userVO := db_user.LoginByToken(userToken)

	a.OK(userVO)
}

func (a UserApi) Logout(c *gin.Context) {
	a.MakeContext(c)
	//token := a.GetToken()
	//userId := a.GetCurrentUserId()
	//db_user.Exit(token, userId)
	a.OK("退出成功")
}

func (a UserApi) UpdateUserInfo(c *gin.Context) {
	a.MakeContext(c)
	var vo vo.UserVO
	if a.Bind(&vo) {
		return
	}

	vo.ID = user.GetUserId32(c)
	ret := db_user.UpdateUserInfo(&vo)
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

	switch flag {
	case 1:
		if a.EmptyFailed("请先绑定手机号！", user.GetUserPhoneNumber(c)) {
			return
		}
		log.Info("GetCaptcha", user.GetUserIdStr(c), "---"+user.GetUserName(c)+"---"+"手机验证码:"+captcha)
	case 2:
		if a.EmptyFailed("请先绑定邮箱！", user.GetEmail(c)) {
			return
		}
		log.Info("GetCaptcha", user.GetUserIdStr(c), "---"+user.GetUserName(c)+"---"+"邮箱验证码:"+captcha)
		sendMail(captcha, user.GetEmail(c))
	default:
		break
	}

	cache.SetCaptchaExpire(user.GetUserIdStr(c), flag, captcha)

	a.OK()
}

// GetCaptchaForBind
// 绑定手机号或者邮箱
// flag: 1 手机号 2 邮箱
func (a UserApi) GetCaptchaForBind(c *gin.Context) {
	a.MakeContext(c)
	//place := a.Param("place")
	//flag := a.Param("flag")
	var param vo.Param
	if a.Bind(&param) {
		return
	}

	captcha := utils.CreateCaptcha(6)

	switch param.Flag {
	case 1:
		log.Info("GetCodeForBind", param.Place+"---"+"手机验证码:"+captcha)
	case 2:
		log.Info("GetCodeForBind", param.Place+"---"+"邮箱验证码:"+captcha)
		sendMail(captcha, user.GetEmail(c))
	default:
		break
	}

	cache.SetCaptchaBindExpire(user.GetUserIdStr(c), param.Place, param.Flag, captcha)

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
	var param vo.Param
	if a.Bind(&param) {
		return
	}

	userId := user.GetUserId32(c)
	a.OK(db_user.UpdateSecretInfo(param.Place, param.FlagString(), param.Code, param.Password, userId))
}

// GetCaptchaForForgetPassword
// 忘记密码 获取验证码
// 1 手机号
// 2 邮箱
func (a UserApi) GetCaptchaForForgetPassword(c *gin.Context) {
	a.MakeContext(c)
	place := a.Query("place")
	flag := a.QueryInt("flag")
	param := vo.Param{Place: place, Flag: flag}

	//captcha := "123456"
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

// GetInfo 获取个人信息
func (a UserApi) GetInfo(c *gin.Context) {
	s := service.User{}
	if a.MakeContextChain(c, &s.Service, nil) == nil {
		return
	}

	userId := user.GetUserId32(c)
	u := blog.User{ID: userId}

	if a.IsError(s.GetUser(&u)) {
		return
	}

	r := blog.Role{ID: u.RoleId}
	if a.IsError(s.GetRole(&r)) {
		return
	}

	var mp = make(map[string]interface{})
	mp["user"] = &u
	mp["role"] = &r
	mp["code"] = 200
	a.OK(mp)
}

func (a UserApi) GetInfo0(c *gin.Context) {
	//req := dto.SysUserById{}
	//s := service.SysUser{}
	//r := service.SysRole{}
	//err := a.MakeContext(c).
	//	MakeOrm().
	//	MakeService(&r.Service).
	//	MakeService(&s.Service).
	//	Errors
	//
	//p := actions.GetPermissionFromContext(c)
	////var roles = make([]string, 1)
	//rolesName := user.GetRoleName(c)
	//var permissions = make([]string, 1)
	//permissions[0] = "*:*:*"
	//var buttons = make([]string, 1)
	//buttons[0] = "*:*:*"
	//
	//var mp = make(map[string]interface{})
	//mp["roles"] = roles
	//if user.GetRoleName(c) == "admin" || user.GetRoleName(c) == "系统管理员" {
	//	mp["permissions"] = permissions
	//	mp["buttons"] = buttons
	//} else {
	//	list, _ := r.GetById(user.GetRoleId(c))
	//	mp["permissions"] = list
	//	mp["buttons"] = list
	//}
	//sysUser := models.SysUser{}
	//req.ID = user.GetUserId(c)
	//err = s.Get(&req, p, &sysUser)
	//if err != nil {
	//	a.Error(http.StatusUnauthorized, err, "登录失败")
	//	return
	//}
	//mp["introduction"] = " am a super administrator"
	//mp["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	//if sysUser.Avatar != "" {
	//	mp["avatar"] = sysUser.Avatar
	//}
	//mp["userName"] = sysUser.NickName
	//mp["userId"] = sysUser.UserId
	//mp["deptId"] = sysUser.DeptId
	//mp["name"] = sysUser.NickName
	//mp["code"] = 200
	//a.OK(mp)
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
