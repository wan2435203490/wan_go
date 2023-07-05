package vo

import r "wan_go/pkg/common/response"

type Page[T any] struct {
	r.Pagination `json:",inline"`
	Records      []T `json:"records" form:"records"`
}

//糟糕的泛型
//func NewPage(pagination *r.Pagination) *PageUser[T] {
//
//}

func (page *Page[T]) Set(pagination *r.Pagination) {
	page.Current = pagination.Current
	page.Size = pagination.Size
}
