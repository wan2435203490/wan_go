package compression

import (
	"fmt"
	"testing"
)

func TestGzipEncode(t *testing.T) {

	src := "https://aqqmusic.tc.qq.com/amobile.music.tc.qq.com/C400001PXEWy1CdkYc.m4a?guid=undefined&v.com/amobile.music.tc.qq.com/C400001PXEWy1CdkYc.m4a?guid=undefined&vkey=3F3BFD7C68926F5946B7F5A628044382B8CC22D1863B7A946792B04109B860586.com/amobile.music.tc.qq.com/C400001PXEWy1CdkYc.m4a?guid=undefined&vkey=3F3BFD7C68926F5946B7F5A628044382B8CC22D1863B7A946792B04109B860586.com/amobile.music.tc.qq.com/C400001PXEWy1CdkYc.m4a?guid=undefined&vkey=3F3BFD7C68926F5946B7F5A628044382B8CC22D1863B7A946792B04109B860586.com/amobile.music.tc.qq.com/C400001PXEWy1CdkYc.m4a?guid=undefined&vkey=3F3BFD7C68926F5946B7F5A628044382B8CC22D1863B7A946792B04109B860586.com/amobile.music.tc.qq.com/C400001PXEWy1CdkYc.m4a?guid=undefined&vkey=3F3BFD7C68926F5946B7F5A628044382B8CC22D1863B7A946792B04109B860586key=3F3BFD7C68926F5946B7F5A628044382B8CC22D1863B7A946792B04109B860586E94ECE8B3AF056BAD590C18114E0B82E1294F931AEE0E18&uin=&fromtag=123032"
	got, _ := MarshalJsonAndGzip([]byte(src))
	fmt.Println(len(got), got)

	bs := make([]byte, 0)
	_ = UnmarshalDataFromJsonWithGzip(got, &bs)
	fmt.Println(len(bs), string(bs))
}
