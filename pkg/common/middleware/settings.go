package middleware

type UrlInfo struct {
	Url    string
	Method string
}

// CasbinExclude casbin 排除的路由列表
var CasbinExclude = []UrlInfo{
	{Url: "/test", Method: "GET"},
}
