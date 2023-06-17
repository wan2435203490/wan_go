package api

import (
	"log"
	"net/http"
	"wan_go/pkg/utils"
)

const (
	SuccessStatus = 0
	ErrorStatus   = -1
)

type Response struct {
	Code    int    `json:"code,omitempty"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

type response struct {
	Response
	Data any `json:"data"`
}

// Pagination 分页
// Current == Offset
// Size == Limit
// Desc 、 Column  Order By
type Pagination struct {
	Current int  `json:"current,omitempty"`
	Size    int  `json:"size,omitempty"`
	Total   int  `json:"total,omitempty"`
	Desc    bool `json:"desc,omitempty"`
	//排序的column 需要多个排序的话就将Desc和Column抽象出来 默认按主键排序
	Column string `json:"column,omitempty"`
}

func (pagination *Pagination) Order() string {

	order := utils.IfThen(utils.IsEmpty(pagination.Column), "id", pagination.Column).(string)

	if pagination.Desc {
		order += " desc"
	}

	return order
}

type CodeMsg struct {
	Code   int    `json:"-"`
	ErrMsg string `json:"-"`
}

var (
	PARAMETER_ERROR = CodeMsg{400, "参数异常！"}
	NOT_LOGIN       = CodeMsg{300, "未登录，请登录后再进行操作！"}
	LOGIN_EXPIRED   = CodeMsg{300, "登录已过期，请重新登录！"}
	SYSTEM_REPAIR   = CodeMsg{301, "系统维护中，敬请期待！"}
	FAIL            = CodeMsg{500, "服务异常！"}
)

func (a *Api) Done(err error) {
	if !a.IsError(err) {
		a.OK()
	}
}

func (a *Api) DoneApiErr(apiErr *CodeMsg) {
	if len(apiErr.ErrMsg) == 0 {
		a.OK()
	}
	a.ErrorInternal(apiErr.ErrMsg)
}

func (a *Api) OK(data ...any) {

	//todo 校验errmsg
	//if reflect.TypeOf(data)

	res := &response{}
	res.Message = http.StatusText(http.StatusOK)
	res.Status = SuccessStatus
	res.Code = http.StatusOK
	res.Data = data

	a.Context.JSON(http.StatusOK, res)
	log.Printf("%#v\n", data)
}
func (a *Api) CodeError(msg CodeMsg) {
	a.Error(msg.Code, msg.ErrMsg)
}

func (a *Api) Error(httpStatus int, msg string) {
	res := &response{}
	res.Message = utils.IfThen(len(msg) > 0, msg, http.StatusText(httpStatus)).(string)
	res.Status = ErrorStatus
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
