package blog

import "wan_go/sdk/api"

type BaseRequestVO[T any] struct {
	api.Pagination
	api.CodeMsg
	Records        []T    `json:"records,omitempty"`
	Order          string `json:"order,omitempty"`
	Source         int32  `json:"source,omitempty"`
	CommentType    string `json:"commentType,omitempty"`
	FloorCommentId int    `json:"floorCommentId,omitempty"`
	SearchKey      string `json:"searchKey,omitempty"`
	// 是否推荐[0:否，1:是]
	RecommendStatus bool   `json:"recommendStatus,omitempty"`
	SortId          int32  `json:"sortId,omitempty"`
	LabelId         int32  `json:"labelId,omitempty"`
	UserStatus      bool   `json:"userStatus,omitempty"`
	UserType        int    `json:"userType,omitempty"`
	UserId          int32  `json:"userId,omitempty"`
	ResourceType    string `json:"resourceType,omitempty"`
	Status          bool   `json:"status,omitempty"`
	Classify        string `json:"classify,omitempty"`
}

func (vo *BaseRequestVO[T]) SetRecords(Records *[]T) {
	vo.Total = len(*Records)
	vo.Records = *Records
}
