package qq

type Sort struct {
	Id   int
	Name string
}

var (
	MoRen   = Sort{Id: 1, Name: "默认"}
	ZuiXin  = Sort{Id: 2, Name: "最新"}
	ZuiRe   = Sort{Id: 3, Name: "最热"}
	PingFen = Sort{Id: 4, Name: "评分"}
)
