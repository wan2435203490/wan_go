package test

import (
	"fmt"
	"regexp"
	"testing"
)

func TestRegister(t *testing.T) {

	const REGEX = `^1[3-9][0-9]\\d{8}`
	userName := "17344030701"
	if matched, err := regexp.MatchString(REGEX, userName); matched || err != nil {
		fmt.Println("error")
	}

}

func TestCheckMobile(t *testing.T) {
	phone := "17asasc多大30701"
	// 匹配规则// ^1第一位为一// [345789]{1} 后接一位345789 的数字// \\d \d的转义 表示数字 {9} 接9位
	//$ 结束符
	regRuler := "^1[345789]{1}\\d{9}$"
	//正则调用规则
	reg := regexp.MustCompile(regRuler)
	// 返回 MatchString 是否匹配
	fmt.Println(reg.MatchString(phone))
}
