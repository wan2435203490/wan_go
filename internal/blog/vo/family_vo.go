package vo

import (
	"time"
	"wan_go/pkg/common/db/mysql/blog"
)

//todo time format

type FamilyVO struct {
	ID             int32  `json:"id"`
	UserId         int32  `json:"userId"`
	BgCover        string `json:"bgCover" vd:"@:len($)>0; msg:'背景封面不能为空'"`
	ManCover       string `json:"manCover" vd:"@:len($)>0; msg:'男生头像不能为空'"`
	WomanCover     string `json:"womanCover" vd:"@:len($)>0; msg:'女生头像不能为空'"`
	ManName        string `json:"manName" vd:"@:len($)>0; msg:'男生昵称不能为空'"`
	WomanName      string `json:"womanName" vd:"@:len($)>0; msg:'女生昵称不能为空'"`
	Timing         string `json:"timing"`
	CountdownTitle string `json:"countdownTitle"`
	Status         bool   `json:"status"`
	CountdownTime  string `json:"countdownTime"`
	FamilyInfo     string `json:"familyInfo"`
	//点赞数
	LikeCount int32     `json:"likeCount"`
	CreatedAt time.Time `json:"createdAt"`
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
