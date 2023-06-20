/*
 * @Author: lwnmengjing
 * @Date: 2021/6/9 10:39 上午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/6/9 10:39 上午
 */

package api

import (
	"fmt"
	"wan_go/core/tools/language"
	"wan_go/sdk/pkg"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
)

func (a *Api) Translate(form, to interface{}) {
	pkg.Translate(form, to)
}

// getAcceptLanguage 获取当前语言
func (a *Api) getAcceptLanguage() string {
	languages := language.ParseAcceptLanguage(a.Context.GetHeader("Accept-Language"), nil)
	if len(languages) == 0 {
		return DefaultLanguage
	}
	return languages[0]
}

// transInit local 通常取决于 http 请求头的 'Accept-Language'
func transInit(local string) (trans ut.Translator, err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() //chinese
		enT := en.New() //english
		uni := ut.New(enT, zhT, enT)

		var o bool
		//register translate
		// 注册翻译器
		switch local {
		case "zh", "zh-CN":
			trans, o = uni.GetTranslator("zh")
			if !o {
				return nil, fmt.Errorf("uni.GetTranslator(%s) failed", "zh")
			}
			err = chTranslations.RegisterDefaultTranslations(v, trans)
		default:
			trans, o = uni.GetTranslator("en")
			if !o {
				return nil, fmt.Errorf("uni.GetTranslator(%s) failed", "en")
			}
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}
