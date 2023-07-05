package blog

import (
	r "wan_go/pkg/common/response"
)

type BaseRequestVO[T any] struct {
	r.Pagination   `json:",inline"`
	r.CodeMsg      `json:"-"`
	Records        []T    `json:"records" form:"records"`
	Order          string `json:"order" form:"order"`
	Source         int32  `json:"source" form:"source"`
	CommentType    string `json:"commentType" form:"commentType"`
	FloorCommentId int    `json:"floorCommentId" form:"floorCommentId"`
	SearchKey      string `json:"searchKey" form:"searchKey"`
	// 是否推荐[0:否，1:是]
	RecommendStatus bool   `json:"recommendStatus" form:"recommendStatus"`
	SortId          int32  `json:"sortId" form:"sortId"`
	LabelId         int32  `json:"labelId" form:"labelId"`
	UserStatus      bool   `json:"userStatus" form:"userStatus"`
	UserType        int    `json:"userType" form:"userType"`
	UserId          int32  `json:"userId" form:"userId"`
	ResourceType    string `json:"resourceType" form:"resourceType"`
	Status          bool   `json:"status" form:"status"`
	Classify        string `json:"classify" form:"classify"`
}

func (vo *BaseRequestVO[T]) SetRecords(Records *[]T) {
	vo.Total = len(*Records)
	vo.Records = *Records
}
