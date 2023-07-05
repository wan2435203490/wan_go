package captcha

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestDriverDigitFunc(t *testing.T) {
	bs := []byte("http://music.163.com/song/media/outer/url?id=2025227742")
	des := make([]byte, base64.URLEncoding.EncodedLen(len(bs)))
	ret := base64.URLEncoding.EncodeToString(bs)
	base64.URLEncoding.Encode(des, bs)
	fmt.Println(ret)
	decodeString, _ := base64.URLEncoding.DecodeString(ret)
	fmt.Println(string(decodeString))
	fmt.Println(string(des))
}
