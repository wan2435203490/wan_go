package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"wan_go/pkg/common/config"
)

func TestLoginAesEncrypt(t *testing.T) {
	//123456 want APQHcSKDlE+dVPS7qoPxig==
	fmt.Println(len(AesEncrypt("0fjkj00fjkjhkvb0fjkjhkvb0fjkjhkvb0fjkjhkvbfjkjhkvbhkvb", config.Config.Blog.CrypotJSKey)))
}
func TestAesEncryptCrypotJsKey(t *testing.T) {
	str := AesEncryptCrypotJsKey("qwq")
	fmt.Println(str)

	str = AesDecryptCrypotJsKey(str)
	fmt.Println(str)
}

func TestEnc(t *testing.T) {
	err := bcrypt.CompareHashAndPassword([]byte("$2a$10$WPCDDh8D6eH5AkFq8WaVceF27CfsVTirv.sIJ5nb.z4DdjCGqqzyG"), []byte("qwq"))
	if err != nil {
		println(err)
	}
	//var arr []string
	//for i := 0; i < 10; i++ {
	//	password := "qwq"
	//	fromPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//	fmt.Println(string(fromPassword))
	//	arr = append(arr, string(fromPassword))
	//}
	//
	//for i, v := range arr {
	//	err := bcrypt.CompareHashAndPassword([]byte(v), []byte("qwq"))
	//	if err != nil {
	//		fmt.Println(i, ":", err)
	//	}
	//}

}

func TestAesSimpleEncrypt(t *testing.T) {
	data := "Hello World!"
	keys := []string{
		"1234",
		"16bit secret key",
		"16bit secret key1234567",
		"16bit secret key12345678",
		"16bit secret key16bit secret ke",
		"16bit secret key16bit secret key",
		"16bit secret key16bit secret key1",
	}
	res := []string{
		"NHlpzbcTvOj686VaF7fU7g==",
		"PuMhKY8ZFLnDAwlQ7v/2SQ==",
		"ZG9JUBvEXrXwSS2RIHvpog==",
		"pbvDuBOV3tJrlPV0xdmbKQ==",
		"uAeg71zBzFeUfEMHJqCSxw==",
		"j9SbFFEEFX4dT9VaDAzsCg==",
		"j9SbFFEEFX4dT9VaDAzsCg==",
	}
	for i, key := range keys {
		if got := AesEncrypt(data, key); got != res[i] {
			t.Errorf("key = %s AesEncrypt() = %s, want %v", key, got, res[i])
		}
	}
}

func TestAesSimpleDecrypt(t *testing.T) {
	data := "Hello World!"
	keys := []string{
		"1234",
		"16bit secret key",
		"16bit secret key1234567",
		"16bit secret key12345678",
		"16bit secret key16bit secret ke",
		"16bit secret key16bit secret key",
		"16bit secret key16bit secret key1",
	}
	res := []string{
		"NHlpzbcTvOj686VaF7fU7g==",
		"PuMhKY8ZFLnDAwlQ7v/2SQ==",
		"ZG9JUBvEXrXwSS2RIHvpog==",
		"pbvDuBOV3tJrlPV0xdmbKQ==",
		"uAeg71zBzFeUfEMHJqCSxw==",
		"j9SbFFEEFX4dT9VaDAzsCg==",
		"j9SbFFEEFX4dT9VaDAzsCg==",
	}
	for i, key := range keys {
		if got := AesDecrypt(res[i], key); got != data {
			t.Errorf("key = %s AesEncrypt() = %s, want %v", key, got, data)
		}
	}
}

func TestGenIVFromKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		wantIv string
	}{
		{
			name: "test",
			args: args{
				key: "16bit secret key",
			},
			wantIv: "ba79295cdabd3a86",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIv := GenIVFromKey(tt.args.key); gotIv != tt.wantIv {
				t.Errorf("GenIVFromKey() = %v, want %v", gotIv, tt.wantIv)
			}
		})
	}
}

func TestAesEncrypt(t *testing.T) {
	type args struct {
		data        string
		key         string
		iv          string
		paddingMode PaddingMode
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				data:        "123456",
				key:         config.Config.Blog.CrypotJSKey,
				iv:          GenIVFromKey(config.Config.Blog.CrypotJSKey),
				paddingMode: PKCS7,
			},
			want: "APQHcSKDlE+dVPS7qoPxig==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AesCBCEncrypt(tt.args.data, tt.args.key, tt.args.iv, tt.args.paddingMode); got != tt.want {
				t.Errorf("AesEncrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAesDecrypt(t *testing.T) {
	type args struct {
		data        string
		key         string
		iv          string
		paddingMode PaddingMode
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				data:        "APQHcSKDlE+dVPS7qoPxig==",
				key:         config.Config.Blog.CrypotJSKey,
				iv:          GenIVFromKey(config.Config.Blog.CrypotJSKey),
				paddingMode: PKCS7,
			},
			want: "123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AesCBCDecrypt(tt.args.data, tt.args.key, tt.args.iv, tt.args.paddingMode); got != tt.want {
				t.Errorf("AesDecrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
