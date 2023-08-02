package dto

type Login struct {
	UserName string `json:"username" vd:"@:len($)>0; msg:'username不能为空'"`
	Password string `json:"password" vd:"@:len($)>0; msg:'password不能为空'"`
}
