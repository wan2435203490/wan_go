package dto

import (
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/models"
	r "wan_go/pkg/common/response"
)

// SaveSortReq update and insert
type SaveSortReq struct {
	ID              int32  `uri:"id"`
	SortName        string `json:"sortName" vd:"@:len($)>0; msg:'分类名称不能为空！'"`
	SortDescription string `json:"sortDescription" vd:"@:len($)>0; msg:'分类描述不能为空！'"`
	SortType        int8   `json:"sortType" vd:"$>-1; msg:'分类类型不能为空！'"` //blog_const.SORT_TYPE_BAR
	Priority        int32  `json:"priority" vd:"!((SortType)$==0 && $<=0); msg:'导航栏分类优先级不能为空！'"`
	models.ControlBy
}

func (from *SaveSortReq) CopyTo(to *blog.Sort) {
	if from.ID != 0 {
		to.ID = from.ID
	}
	to.SortName = from.SortName
	to.SortDescription = from.SortDescription
	to.SortType = from.SortType
	to.Priority = from.Priority
}

func (s *SaveSortReq) GetId() interface{} {
	return s.ID
}

type DelSortReq struct {
	Ids []int `json:"ids"`
}

func (s *DelSortReq) GetId() interface{} {
	return s.Ids
}

type PageSortReq struct {
	*r.Pagination   `json:",inline"`
	ID              int32  `uri:"id" search:"type:eq;column:id;table:sort" comment:"id"`
	SortName        string `form:"sortName" search:"type:like;column:sort_name;table:sort" comment:"分类名称"`
	SortDescription string `form:"sortDescription" search:"type:like;column:sort_description;table:sort" comment:"分类描述"`
	SortType        int8   `form:"sortType" search:"type:eq;column:sort_type;table:sort"` //blog_const.SORT_TYPE_BAR
	Priority        int32  `form:"priority" search:"type:eq;column:sort_priority;table:sort"`
}

func (s *PageSortReq) GetNeedSearch() interface{} {
	return *s
}
