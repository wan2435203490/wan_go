package blog

import "time"

//todo time format

type FamilyVO struct {
	ID             int    `json:"id,omitempty"`
	UserId         int    `json:"userId,omitempty"`
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
	LikeCount  int       `json:"likeCount,omitempty"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}
