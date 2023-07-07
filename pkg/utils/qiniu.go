package utils

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
	"os"
	"wan_go/pkg/common/config"
)

//https://developer.qiniu.com/kodo/1238/go

//1.业务服务器颁发 上传凭证给客户端（终端用户）
//2.客户端凭借上传凭证上传文件到七牛
//3.在七牛获得完整数据后，发起一个 HTTP 请求回调到业务服务器
//4.业务服务器保存相关信息，并返回一些信息给七牛
//5.七牛原封不动地将这些信息转发给客户端（终端用户）

func UploadToQiNiu(file *os.File) {
	putPlicy := storage.PutPolicy{
		Scope:   config.Config.Qiniu.Bucket,
		Expires: 3600, //default could set config
	}
	mac := GetMac()

	// 获取上传凭证
	upToken := putPlicy.UploadToken(mac)

	// 配置参数
	cfg := storage.Config{
		//Zone:          &storage.ZoneHuanan, // 华南区
		UseCdnDomains: true,
		UseHTTPS:      true, // 非https
	}
	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}        // 上传后返回的结果
	putExtra := storage.PutExtra{} // 额外参数

	// 上传 自定义key，可以指定上传目录及文件名和后缀，
	key := "test/" + file.Name() // 上传路径，如果当前目录中已存在相同文件，则返回上传失败错误
	finfo, _ := os.Stat(file.Name())
	err := formUploader.Put(context.Background(), &ret, upToken, key, file, finfo.Size(), &putExtra)

	// 以默认key方式上传
	// err = formUploader.PutWithoutKey(context.Background(), &ret, upToken, src, fileSize, &putExtra)

	// 自定义key，上传指定路径的文件
	// localFilePath = "./aa.jpg"
	// err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFilePath, &putExtra)

	// 默认key，上传指定路径的文件
	// localFilePath = "./aa.jpg"
	// err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFilePath, &putExtra)

	if err != nil {
		fmt.Println(err.Error())
	}

	url := config.Config.Qiniu.Url + ret.Key // 返回上传后的文件访问路径

	fmt.Println(url)
}

// UploadImgToQiNiu 上传图片到七牛云，然后返回状态和图片的url
func UploadImgToQiNiu(file *multipart.FileHeader) (int, string) {

	src, err := file.Open()
	if err != nil {
		return 10011, err.Error()
	}
	defer src.Close()

	putPlicy := storage.PutPolicy{
		Scope:   config.Config.Qiniu.Bucket,
		Expires: 3600, //default could set config
	}
	mac := qbox.NewMac(config.Config.Qiniu.AccessKey, config.Config.Qiniu.SecretKey)

	// 获取上传凭证
	upToken := putPlicy.UploadToken(mac)

	// 配置参数
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan, // 华南区
		UseCdnDomains: true,
		UseHTTPS:      true, // 非https
	}
	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}        // 上传后返回的结果
	putExtra := storage.PutExtra{} // 额外参数

	// 上传 自定义key，可以指定上传目录及文件名和后缀，
	key := "image/" + file.Filename // 上传路径，如果当前目录中已存在相同文件，则返回上传失败错误
	err = formUploader.Put(context.Background(), &ret, upToken, key, src, file.Size, &putExtra)

	// 以默认key方式上传
	// err = formUploader.PutWithoutKey(context.Background(), &ret, upToken, src, fileSize, &putExtra)

	// 自定义key，上传指定路径的文件
	// localFilePath = "./aa.jpg"
	// err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFilePath, &putExtra)

	// 默认key，上传指定路径的文件
	// localFilePath = "./aa.jpg"
	// err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFilePath, &putExtra)

	if err != nil {
		code := 501
		return code, err.Error()
	}

	url := config.Config.Qiniu.Url + ret.Key // 返回上传后的文件访问路径
	return 0, url
}

func GetMac() *qbox.Mac {
	return qbox.NewMac(config.Config.Qiniu.AccessKey, config.Config.Qiniu.SecretKey)
}

func Scope(key string) string {
	if key == "" {
		return config.Config.Qiniu.Bucket
	}
	return fmt.Sprintf("%s:%s", config.Config.Qiniu.Bucket, key)
}

func BucketManager() *storage.BucketManager {
	mac := GetMac()
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: true,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Region=&storage.ZoneHuabei
	return storage.NewBucketManager(mac, &cfg)
}

// GetQiniuToken 获取上传凭证
func GetQiniuToken(key string) string {
	putPolicy := storage.PutPolicy{
		Scope: Scope(key),
		//Expires: 3600,//default
		//七牛云回复格式
		//ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
		//上传到七牛云后需要回调服务器时需要callback
		//CallbackURL:      "http://api.example.com/qiniu/upload/callback",
		//CallbackBody:     `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
		//CallbackBodyType: "application/json",
	}
	mac := GetMac()

	// 获取上传凭证
	upToken := putPolicy.UploadToken(mac)
	return upToken
}

func DeleteQiniuFile(field string) {

	err := BucketManager().Delete(config.Config.Qiniu.Bucket, field)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//	BatchDelete fields := []string{
//		"github1.png",
//		"github2.png",
//	}
func BatchDelete(fields []string) {
	//每个batch的操作数量不可以超过1000个，如果总数量超过1000，需要分批发送
	deleteOps := make([]string, 0, len(fields))
	for _, key := range fields {
		deleteOps = append(deleteOps, storage.URIDelete(config.Config.Qiniu.Bucket, key))
	}

	rets, err := BucketManager().Batch(deleteOps)
	if len(rets) == 0 {
		// 处理错误
		if e, ok := err.(*storage.ErrorInfo); ok {
			fmt.Printf("batch error, code:%d", e.Code)
		} else {
			fmt.Printf("batch error, %s", err)
		}
		return
	}

	// 返回 rets，先判断 rets 是否
	for _, ret := range rets {
		// 200 为成功
		if ret.Code == 200 {
			fmt.Printf("%d\n", ret.Code)
		} else {
			fmt.Printf("%s\n", ret.Data.Error)
		}
	}

}

func BatchGetFiles(keys []string) map[string]map[string]string {
	//每个batch的操作数量不可以超过1000个，如果总数量超过1000，需要分批发送
	getOps := make([]string, 0, len(keys))
	for _, key := range keys {
		getOps = append(getOps, storage.URIStat(config.Config.Qiniu.Bucket, key))
	}

	rets, err := BucketManager().Batch(getOps)
	if len(rets) == 0 {
		// 处理错误
		if _, ok := err.(*storage.ErrorInfo); ok {
			//logs.NewWarn("BatchGetFiles", "batch error, code:%d", e.Captcha)
		} else {
			//logs.Info("BatchGetFiles", "batch error, %s", err)
		}
		return nil
	}

	result := make(map[string]map[string]string, 16)
	for i, key := range keys {
		ret := rets[i]
		if ret.Code == 200 {
			info := make(map[string]string, 16)
			info["size"] = int64ToString(ret.Data.Fsize)
			info["mimeType"] = ret.Data.MimeType
			result[key] = info
		} else {
			//logs.NewWarn("BatchGetFiles", "%s\n", ret.Data.Error)
		}
	}

	return result
}

//
//func uptoken(bucketName string) string {
//	putPolicy := rs.PutPolicy{
//		Scope: bucketName,
//		//CallbackUrl: callbackUrl,
//		//CallbackBody:callbackBody,
//		//ReturnUrl:   returnUrl,
//		//ReturnBody:  returnBody,
//		//AsyncOps:    asyncOps,
//		//EndUser:     endUser,
//		//Expires:     expires,
//	}
//	return putPolicy.Token(nil)
//} // 文件上传
//func TestUploadFile(t *testing.T) {
//	var ak = core.QiNiuAK
//	var sk = core.QiNiuSk
//	var bucket = core.QiNiuBucket
//	var url = core.QiuNiuUrl
//
//	src, err := os.ReadFile("./img/meinv.jpeg")
//	if err != nil {
//		t.Fatal(err)
//	}
//	fileSize := len(src)
//
//	putPolicy := storage.PutPolicy{
//		Scope: bucket,
//	}
//	mac := qbox.NewMac(ak, sk)
//	upToken := putPolicy.UploadToken(mac)
//	cfg := storage.Config{
//		Zone:          &storage.ZoneHuanan,
//		UseCdnDomains: false,
//		UseHTTPS:      false,
//	}
//	putExtra := storage.PutExtra{}
//	formUploader := storage.NewFormUploader(&cfg)
//	ret := storage.PutRet{}
//	key := "go-cloud-storage/meinv.jpeg"
//	err = formUploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(src), int64(fileSize), &putExtra)
//	if err != nil {
//		t.Fatal(err)
//	}
//	url2 := url + ret.Key
//	fmt.Println(ret)
//	fmt.Println(url2)
//}
//
//// 分片上传
//func TestUploadChunkFile(t *testing.T) {
//	var ak = core.QiNiuAK
//	var sk = core.QiNiuSk
//	var bucket = core.QiNiuBucket
//	var url = core.QiuNiuUrl
//
//	putPolicy := storage.PutPolicy{
//		Scope: bucket,
//	}
//	mac := qbox.NewMac(ak, sk)
//	upToken := putPolicy.UploadToken(mac)
//	cfg := storage.Config{
//		Zone:          &storage.ZoneHuanan,
//		UseCdnDomains: false,
//		UseHTTPS:      false,
//	}
//	resumeUploaderV2 := storage.NewResumeUploaderV2(&cfg)
//	upHost, err := resumeUploaderV2.UpHost(ak, bucket)
//	if err != nil {
//		t.Fatal(err)
//	}
//	key := "go-cloud-storage/lala.mp4"
//	// 初始化分块上传
//	initPartsRet := storage.InitPartsRet{}
//	err = resumeUploaderV2.InitParts(context.TODO(), upToken, upHost, bucket, key, true, &initPartsRet)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fileInfo, err := os.Open("./music/lala.mp4")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer fileInfo.Close()
//	fileContent, err := ioutil.ReadAll(fileInfo)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fileLen := len(fileContent)
//	chunkSize2 := 2 * 1024 * 1024
//
//	num := fileLen / chunkSize2
//	if fileLen%chunkSize2 > 0 {
//		num++
//	}
//
//	// 分块上传
//	var uploadPartInfos []storage.UploadPartInfo
//	for i := 1; i <= num; i++ {
//		partNumber := int64(i)
//		fmt.Printf("开始上传第%v片数据", partNumber)
//
//		var partContentBytes []byte
//		endSize := i * chunkSize2
//		if endSize > fileLen {
//			endSize = fileLen
//		}
//		partContentBytes = fileContent[(i-1)*chunkSize2 : endSize]
//		partContentMd5 := Md5(string(partContentBytes))
//		uploadPartsRet := storage.UploadPartsRet{}
//		err = resumeUploaderV2.UploadParts(context.TODO(), upToken, upHost, bucket, key, true,
//			initPartsRet.UploadID, partNumber, partContentMd5, &uploadPartsRet, bytes.NewReader(partContentBytes),
//			len(partContentBytes))
//		if err != nil {
//			t.Fatal(err)
//		}
//		uploadPartInfos = append(uploadPartInfos, storage.UploadPartInfo{
//			Etag:       uploadPartsRet.Etag,
//			PartNumber: partNumber,
//		})
//		fmt.Printf("结束上传第%d片数据\n", partNumber)
//	}
//
//	// 完成上传
//	rPutExtra := storage.RputV2Extra{Progresses: uploadPartInfos}
//	comletePartRet := storage.PutRet{}
//	err = resumeUploaderV2.CompleteParts(context.TODO(), upToken, upHost, &comletePartRet, bucket, key,
//		true, initPartsRet.UploadID, &rPutExtra)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	url2 := url + comletePartRet.Key
//	fmt.Println(comletePartRet.Hash)
//	fmt.Println(url2)
//}
//
//// 断点续传
//func TestResumeUploadFile(t *testing.T) {
//	ak := core.QiNiuAK
//	sk := core.QiNiuSk
//	localFile := "./music/abc.mp4"
//	bucket := core.QiNiuBucket
//	key := "go-cloud-storage/abc.mp4"
//	url := core.QiuNiuUrl
//	putPolicy := storage.PutPolicy{
//		Scope: bucket,
//	}
//	mac := qbox.NewMac(ak, sk)
//	upToken := putPolicy.UploadToken(mac)
//	cfg := storage.Config{
//		Zone:          &storage.ZoneHuanan,
//		UseCdnDomains: false,
//		UseHTTPS:      false,
//	}
//	resumeUploaderV2 := storage.NewResumeUploaderV2(&cfg)
//	ret := storage.PutRet{}
//	recorder, err := storage.NewFileRecorder(os.TempDir())
//	if err != nil {
//		t.Fatal(err)
//	}
//	putExtra := storage.RputV2Extra{
//		Recorder: recorder,
//	}
//	err = resumeUploaderV2.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	url2 := url + ret.Key
//	fmt.Println(ret)
//	fmt.Println(url2)
//}
