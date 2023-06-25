package api

import (
	json2 "encoding/json"
	"log"
	"net/http"
	"reflect"
	r "wan_go/pkg/common/response"
	"wan_go/pkg/utils"
)

func (a *Api) Done(err error) {
	if !a.IsError(err) {
		a.OK()
	}
}

func (a *Api) DoneApiErr(apiErr *r.CodeMsg) {
	if apiErr == nil || len(apiErr.Msg) == 0 {
		a.OK()
		return
	}
	a.ErrorInternal(apiErr.Msg)
}

func (a *Api) OK(data ...any) {

	//todo 校验errmsg
	if data != nil && data[0] != nil {
		value := reflect.ValueOf(data[0])
		if value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface {
			value = value.Elem()
		}
		if value.Kind() == reflect.Struct {
			errMsgValue := value.FieldByName("Msg")
			if errMsgValue.IsValid() && !errMsgValue.IsZero() {

				codeValue := value.FieldByName("Code")
				var code int
				if codeValue.IsZero() {
					code = 400
				} else {
					code = int(codeValue.Int())
				}
				a.Error(code, errMsgValue.String())
				return
			}
		}
	}

	res := &r.Response{}
	res.Message = http.StatusText(http.StatusOK)
	res.Status = r.SuccessStatus
	res.Code = http.StatusOK
	if data != nil {
		res.Data = data[0]
	}

	a.Context.JSON(http.StatusOK, res)
	//log.Println(a.Context.Request.URL, &a.Context, &res, "\n")
	rr, _ := json2.Marshal(res)
	log.Println(a.Context.Request.URL, ":", string(rr), "\n")
}
func (a *Api) CodeError(msg r.CodeMsg) {
	a.Error(msg.Code, msg.Msg)
}

func (a *Api) Error(httpStatus int, msg string) {
	res := &r.Response{}
	res.Message = utils.IfThen(len(msg) > 0, msg, http.StatusText(httpStatus)).(string)
	res.Status = r.ErrorStatus
	res.Code = httpStatus

	a.Context.JSON(httpStatus, res)
	log.Printf("%#v\n", msg)
}

func (a *Api) ErrorInternal(msg string) {
	a.Error(http.StatusInternalServerError, msg)
}

// IsFailed named what?
func (a *Api) IsFailed(cond bool, msg string) bool {

	if cond {
		a.Error(http.StatusInternalServerError, msg)
	}

	return cond
}

// EmptyFailed 校验 strs 有空字符串时 返回 msg
func (a *Api) EmptyFailed(msg string, strs ...string) bool {
	return a.IsFailed(utils.IsEmpty(strs...), msg)
}
