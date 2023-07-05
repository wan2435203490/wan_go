package api

import (
	json2 "encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	r "wan_go/pkg/common/response"
	"wan_go/pkg/utils"
	"wan_go/sdk/pkg/response"
)

func (a Api) Done(err error) {
	if !a.IsError(err) {
		a.OK()
	}
}

func (a Api) DoneApiErr(apiErr *r.CodeMsg) {
	if apiErr == nil || len(apiErr.Msg) == 0 {
		a.OK()
		return
	}
	a.ErrorInternal(apiErr.Msg)
}

func (a Api) OKMsg(data interface{}, msg string) {
	response.OK(a.Context, data, msg)
}

func (a Api) OK(data ...any) {

	//todo 校验errmsg
	//if data != nil && data[0] != nil {
	//	value := reflect.ValueOf(data[0])
	//	if value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface {
	//		value = value.Elem()
	//	}
	//	if value.Kind() == reflect.Struct {
	//		errMsgValue := value.FieldByName("Message")
	//		if errMsgValue.IsValid() && !errMsgValue.IsZero() {
	//
	//			codeValue := value.FieldByName("Captcha")
	//			var code int
	//			if codeValue.IsZero() {
	//				code = 400
	//			} else {
	//				code = int(codeValue.Int())
	//			}
	//			a.ErrorMsg(code, errMsgValue.String())
	//			return
	//		}
	//	}
	//}

	a.OKMsg(data[0], "")

	rr, _ := json2.Marshal(data[0])
	log.Println(a.Context.Request.URL, ":", string(rr), "\n")
}

func (a Api) CodeError(msg r.CodeMsg) {
	a.ErrorMsg(msg.Code, msg.Msg)
}

func (a Api) ErrorMsg(code int, msg string) {

	response.Error(a.Context, code, nil, msg)
}

func (a Api) Error(code int, err error) {
	response.Error(a.Context, code, err, "")

}

func (a Api) ErrorInternal(msg string) {
	a.ErrorMsg(http.StatusInternalServerError, msg)
}

// PageOK 分页数据处理
func (e Api) PageOK(result interface{}, count int, pageIndex int, pageSize int, msg string) {
	response.PageOK(e.Context, result, count, pageIndex, pageSize, msg)
}

// Custom 兼容函数
func (e Api) Custom(data gin.H) {
	response.Custum(e.Context, data)
}

// IsFailed named what?
func (a Api) IsFailed(cond bool, msg string) bool {

	if cond {
		a.ErrorMsg(http.StatusInternalServerError, msg)
	}

	return cond
}

// EmptyFailed 校验 strs 有空字符串时 返回 msg
func (a Api) EmptyFailed(msg string, strs ...string) bool {
	return a.IsFailed(utils.IsEmpty(strs...), msg)
}
