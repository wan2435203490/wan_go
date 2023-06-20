package blog

import (
	"strconv"
)

type Param struct {
	Place    string `json:"place"`
	Flag     int    `json:"flag"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

func (form *Param) FlagString() string {
	return strconv.FormatInt(int64(form.Flag), 10)
}
