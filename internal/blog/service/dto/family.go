package dto

import (
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/common/models"
	r "wan_go/pkg/common/response"
)

// SaveFamilyReq update and insert
type SaveFamilyReq struct {
	//ID                int32  `uri:"id"`
	//SortId            int32  `json:"sortId" vd:"$>0; msg:'分类Id不能为空！'"`
	//FamilyName        string `json:"labelName" vd:"@:len($)>0; msg:'标签名称不能为空！'"`
	//FamilyDescription string `json:"labelDescription" vd:"@:len($)>0; msg:'标签描述不能为空！'"`
	ID               int32  `uri:"id"`
	UserId           int32  `json:"userId,omitempty"`
	BgCover          string `json:"bgCover,omitempty"`
	ManCover         string `json:"manCover,omitempty"`
	WomanCover       string `json:"womanCover,omitempty"`
	ManName          string `json:"manName,omitempty"`
	WomanName        string `json:"womanName,omitempty"`
	Timing           string `json:"timing,omitempty"`
	CountdownTitle   string `json:"countdownTitle,omitempty"`
	CountdownTime    string `json:"countdownTime,omitempty"`
	Status           bool   `json:"status,omitempty"`
	FamilyInfo       string `json:"familyInfo,omitempty"`
	LikeCount        int32  `json:"likeCount,omitempty"`
	models.ControlBy `json:"models.ControlBy"`
}

func (from *SaveFamilyReq) CopyTo(to *blog.Family) {
	if from.ID != 0 {
		to.ID = from.ID
	}
	to.UserId = from.UserId
	to.BgCover = from.BgCover
	to.ManCover = from.ManCover
	to.WomanCover = from.WomanCover
	to.ManName = from.ManName
	to.WomanName = from.WomanName
	to.Timing = from.Timing
	to.CountdownTitle = from.CountdownTitle
	to.CountdownTime = from.CountdownTime
	to.Status = from.Status
	to.FamilyInfo = from.FamilyInfo
	to.LikeCount = from.LikeCount
}

func (s *SaveFamilyReq) GetId() interface{} {
	return s.ID
}

type DelFamilyReq struct {
	Ids []int `json:"ids"`
}

func (s *DelFamilyReq) GetId() interface{} {
	return s.Ids
}

type PageFamilyReq struct {
	*r.Pagination `json:",inline"`
	ID            int32 `uri:"id" search:"type:eq;column:id;table:family" comment:"id"`
	Status        *bool `form:"status" search:"type:eq;column:status;table:family" `
	Size          int   `form:"size"`
}

func (s *PageFamilyReq) GetNeedSearch() interface{} {
	return *s
}

type ChangeFamilyReq struct {
	ID     int32 `uri:"id"`
	Status *bool `json:"status"`
}
