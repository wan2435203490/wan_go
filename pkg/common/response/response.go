package r

import "wan_go/pkg/utils"

const (
	SuccessStatus = 0
	ErrorStatus   = -1
)

type Response struct {
	Code    int    `json:"code,omitempty"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
}

//type response struct {
//	Response
//	Data any `json:"data"`
//}

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
	Code int    `json:"-"`
	Msg  string `json:"-"`
}

var (
	PARAMETER_ERROR = CodeMsg{400, "参数异常！"}
	NOT_LOGIN       = CodeMsg{300, "未登录，请登录后再进行操作！"}
	LOGIN_EXPIRED   = CodeMsg{300, "登录已过期，请重新登录！"}
	SYSTEM_REPAIR   = CodeMsg{301, "系统维护中，敬请期待！"}
	FAIL            = CodeMsg{500, "服务异常！"}
)
