package blog

import "wan_go/sdk/api"

type BaseRequestVO[T any] struct {
	api.Pagination
	Records        []T    `json:"records,omitempty"`
	Order          string `json:"order,omitempty"`
	Source         int    `json:"source,omitempty"`
	CommentType    string `json:"commentType,omitempty"`
	FloorCommentId int    `json:"floorCommentId,omitempty"`
	SearchKey      string `json:"searchKey,omitempty"`
	// 是否推荐[0:否，1:是]
	RecommendStatus bool   `json:"recommendStatus,omitempty"`
	SortId          int    `json:"sortId,omitempty"`
	LabelId         int    `json:"labelId,omitempty"`
	UserStatus      bool   `json:"userStatus,omitempty"`
	UserType        int    `json:"userType,omitempty"`
	UserId          int    `json:"userId,omitempty"`
	ResourceType    string `json:"resourceType,omitempty"`
	Status          bool   `json:"status,omitempty"`
	Classify        string `json:"classify,omitempty"`
}
