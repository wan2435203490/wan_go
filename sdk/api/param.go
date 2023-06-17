package api

import (
	vd "github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"reflect"
	"strconv"
	"wan_go/pkg/utils"
)

func (a *Api) Param(key string) (value string) {

	value = a.Context.Param(key)
	if value == "" {
		a.Context.AbortWithStatusJSON(http.StatusBadRequest, key+" is empty")
	}

	return
}

func (a *Api) StringFailed(parseStr *string, key string) (failed bool) {

	*parseStr = a.Context.Param(key)
	if *parseStr == "" {
		a.Context.AbortWithStatusJSON(http.StatusBadRequest, key+" is empty")
		failed = true
	}

	return
}

func (a *Api) BoolFailed(parseBool *bool, key string) (failed bool) {

	value := a.Context.Param(key)

	aBool, err := strconv.ParseBool(value)

	if value == "" || err != nil {
		a.Context.AbortWithStatusJSON(http.StatusBadRequest, key+" is empty")
		failed = true
	}

	*parseBool = aBool

	return
}

func (a *Api) IntFailed(parseInt *int, key string) (failed bool) {

	value := a.Context.Param(key)

	if value == "" {
		a.Context.AbortWithStatusJSON(http.StatusBadRequest, key+" is empty")
		failed = true
	}

	*parseInt = utils.StringToInt(value)

	return
}

// BindFailed validate && write error json response if error
func (a *Api) BindFailed(d interface{}, bindings ...binding.Binding) bool {
	if d == nil {
		return true
	}

	var err error
	if len(bindings) == 0 {
		bindings = constructor.GetBindingForGin(d)
	}
	for i := range bindings {
		if bindings[i] == nil {
			err = a.Context.ShouldBindUri(d)
		} else {
			err = a.Context.ShouldBindWith(d, bindings[i])
		}
		if err != nil && err.Error() == "EOF" {
			//a.Logger.Warn("request body is not present anymore. ")
			err = nil
			continue
		}
		if err != nil {
			a.AddError(err)
			break
		}
	}

	if err1 := vd.Validate(d); err1 != nil {
		a.AddError(err1)
	}

	if a.Err != nil {
		a.ErrorInternal(a.Err.Error())
		return true
	}

	//todo set status when d is blogVO.BaseRequestVO
	if a.IsAdmin() && reflect.TypeOf(d).String() == "blogVO.BaseRequestVO" {
		//todo test
		v := reflect.ValueOf(d)
		v.FieldByName("Status").Set(reflect.ValueOf(true))
	}

	return false
}
