package antd

import (
	"fmt"

	resp "wan_go/sdk/pkg/response"
)

const (
	Silent       = "0"
	MessageWarn  = "1"
	MessageError = "2"
	Notification = "4"
	Page         = "9"
)

type Response struct {
	Success      bool   `json:"success"`      // if request is success
	ErrorCode    string `json:"errorCode"`    // code for errorType
	ErrorMessage string `json:"errorMessage"` // message display to db_user
	ShowType     string `json:"showType"`     // error display type： 0 silent; 1 message.warn; 2 message.error; 4 notification; 9 page
	TraceId      string `json:"traceId"`      // Convenient for back-end Troubleshooting: unique request ID
	Host         string `json:"host"`         // onvenient for backend Troubleshooting: host of current access server
	Status       string `json:"status"`
}
type response struct {
	Response
	Data interface{} `json:"data"` // response data
}

type Pages struct {
	Response
	Data     interface{} `json:"data"` // response data
	Total    int         `json:"total"`
	Current  int         `json:"current"`
	PageSize int         `json:"pageSize"`
}

type pages struct {
	Pages
	Data interface{} `json:"data"`
}

type lists struct {
	Response
	ListData ListData `json:"data"` // response data
}

type ListData struct {
	List     interface{} `json:"list"` // response data
	Total    int         `json:"total"`
	Current  int         `json:"current"`
	PageSize int         `json:"pageSize"`
}

func (e *response) SetCode(code int32) {
	switch code {
	case 200, 0:
	default:
		e.ErrorCode = fmt.Sprintf("C%d", code)
	}
}

func (e *response) SetTraceID(id string) {
	e.TraceId = id
}

func (e *response) SetMsg(msg string) {
	e.ErrorMessage = msg
}

func (e *response) SetData(data interface{}) {
	e.Data = data
}

func (e *response) SetSuccess(success bool) {
	e.Success = success
}

func (e response) Clone() resp.Responses {
	return &e
}
