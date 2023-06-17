package blog

import (
	"time"
	"wan_go/pkg/common/db/mysql/blog"
)

//todo time format

type FamilyVO struct {
	ID             int32  `json:"id,omitempty"`
	UserId         int32  `json:"userId,omitempty"`
	BgCover        string `json:"bgCover,omitempty" vd:"@:len($)>0; msg:'背景封面不能为空'"`
	ManCover       string `json:"manCover,omitempty" vd:"@:len($)>0; msg:'男生头像不能为空'"`
	WomanCover     string `json:"womanCover,omitempty" vd:"@:len($)>0; msg:'女生头像不能为空'"`
	ManName        string `json:"manName,omitempty" vd:"@:len($)>0; msg:'男生昵称不能为空'"`
	WomanName      string `json:"womanName,omitempty" vd:"@:len($)>0; msg:'女生昵称不能为空'"`
	Timing         string `json:"timing,omitempty"`
	CountdownTitle string `json:"countdownTitle,omitempty"`
	Status         bool   `json:"status,omitempty"`
	CountdownTime  string `json:"countdownTime,omitempty"`
	FamilyInfo     string `json:"familyInfo,omitempty"`
	//点赞数
	LikeCount int32     `json:"likeCount,omitempty"`
	CreatedAt time.Time `json:"createTime"`
	UpdatedAt time.Time `json:"updateTime"`
}

func (to *FamilyVO) Copy(from *blog.Family) {
	to.ID = from.ID
	to.UserId = from.UserId
	to.BgCover = from.BgCover
	to.ManCover = from.ManCover
	to.WomanCover = from.WomanCover
	to.ManName = from.ManName
	to.WomanName = from.WomanName
	to.Timing = from.Timing
	to.CountdownTitle = from.CountdownTitle
	to.Status = from.Status
	to.CountdownTime = from.CountdownTime
	to.LikeCount = from.LikeCount
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
}

func (from *FamilyVO) CopyTo(to *blog.Family) {
	to.ID = from.ID
	to.UserId = from.UserId
	to.BgCover = from.BgCover
	to.ManCover = from.ManCover
	to.WomanCover = from.WomanCover
	to.ManName = from.ManName
	to.WomanName = from.WomanName
	to.Timing = from.Timing
	to.CountdownTitle = from.CountdownTitle
	to.Status = from.Status
	to.CountdownTime = from.CountdownTime
	to.LikeCount = from.LikeCount
	to.CreatedAt = from.CreatedAt
	to.UpdatedAt = from.UpdatedAt
}
