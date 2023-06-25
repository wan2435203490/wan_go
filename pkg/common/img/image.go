package app

import (
	"bytes"
	"fmt"

	//"code.google.com/p/graphics-go/graphics"
	"errors"
	"github.com/chai2010/webp"
	"github.com/gogf/gf/os/gfile"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
	"os"
	"strings"
)

const ImageJpeg int = 1
const ImagePng int = 2
const ImageBmp int = 3
const ImageGif int = 4
const ImageWebp int = 5

type Image struct {
	image.Image
	FilePath  string
	Data      []byte
	ImageType int
	Ext       string
	Width     int
	Height    int
}

// Open /**
func (that *Image) Open(filePath string) (err error) {
	that.FilePath = filePath
	that.Data = gfile.GetBytes(filePath)
	that.Ext = gfile.ExtName(filePath)
	contentType := http.DetectContentType(that.Data[:512])
	if strings.Contains(contentType, "jpeg") {
		that.ImageType = ImageJpeg
	} else if strings.Contains(contentType, "png") {
		that.ImageType = ImagePng
	} else if strings.Contains(contentType, "bmp") {
		that.ImageType = ImageBmp
	} else if strings.Contains(contentType, "gif") {
		that.ImageType = ImageGif
	} else if strings.Contains(contentType, "webp") {
		that.ImageType = ImageWebp
	}
	reader := bytes.NewReader(that.Data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}
	b := img.Bounds()
	that.Width = b.Max.X
	that.Height = b.Max.Y
	that.Image = img
	return nil
}

func (that *Image) Reset() {
	that.Data = nil
}

const DEFAULT_MAX_WIDTH float64 = 120
const DEFAULT_MAX_HEIGHT float64 = 120

// 计算图片缩放后的尺寸 缩放到指定大小
func calculateRatioFit(srcWidth, srcHeight int) (int, int) {
	//ratioW, ratioH := DEFAULT_MAX_WIDTH/float64(srcWidth), DEFAULT_MAX_HEIGHT/float64(srcHeight)
	////ratio := math.Min(DEFAULT_MAX_WIDTH/float64(srcWidth), DEFAULT_MAX_HEIGHT/float64(srcHeight))
	//return int(math.Ceil(float64(srcWidth) * ratioW)), int(math.Ceil(float64(srcHeight) * ratioH))
	return 120, 120
}

// 计算图片缩放后的尺寸 等比例缩放
func (that *Image) calculateRatioFit(srcWidth, srcHeight int, desWidth, desHeight int) (int, int) {
	return desWidth, desHeight
}

func GetThumbPath(path string) string {
	//path := "/Users/wan/Pictures/download/jpeg/1.jpeg"
	index := strings.LastIndex(path, `/`) + 1
	str2 := path[:index] + "thumb/120_120_" + path[index:]
	//fmt.Println(str2)
	return str2
}

// 生成缩略图
func MakeJpegThumbnail(imagePath string) error {

	file, _ := os.Open(imagePath)
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	b := img.Bounds()
	width := b.Max.X
	height := b.Max.Y

	w, h := calculateRatioFit(width, height)

	fmt.Println("width = ", width, " height = ", height)
	fmt.Println("w = ", w, " h = ", h)

	// 调用resize库进行图片缩放
	m := resize.Resize(uint(w), uint(h), img, resize.Lanczos3)

	// 需要保存的文件
	savePath := GetThumbPath(imagePath)
	imgfile, _ := os.Create(savePath)
	defer imgfile.Close()

	// 以JPEG格式保存文件
	err = jpeg.Encode(imgfile, m, &jpeg.Options{Quality: 80})
	if err != nil {
		return err
	}

	return nil
}

func (that *Image) JpegToThumb(width, height, quality int, savePath string) error {
	//var img image.Image
	//var f bytes.Buffer
	f, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	subImg := CropImage(that.Image, 0, width, 0, height)
	if subImg == nil {
		return errors.New("Crop error")
	}

	err = jpeg.Encode(f, subImg, &jpeg.Options{Quality: quality})
	if err != nil {
		return err
	}

	return nil
}

// CropImage will get sub image on position
func CropImage(img image.Image, minX, maxX, minY, maxY int) image.Image {
	// check range of position
	if checkRect(minX, maxX, minY, maxY, img.Bounds()) {
		return nil
	}
	// image in memory is image.YCbCr
	rgbImg, ok := img.(*image.YCbCr)
	if !ok {
		return nil
	}
	return rgbImg.SubImage(image.Rect(minX, minY, maxX, maxY))
}

func checkRect(minX, maxX, minY, maxY int, r image.Rectangle) bool {
	if maxX > r.Dx() || maxY > r.Dy() {
		return false
	}
	return minX < 0 || minX >= maxX || minY < 0 || minY >= maxY
}

const MaxWidth float64 = 600

func fixSize(img1W, img2W int) (new1W, new2W int) {
	var ( //为了方便计算，将两个图片的宽转为 float64
		img1Width, img2Width = float64(img1W), float64(img2W)
		ratio1, ratio2       float64
	)

	minWidth := math.Min(img1Width, img2Width) // 取出两张图片中宽度最小的为基准

	if minWidth > 600 { // 如果最小宽度大于600，那么两张图片都需要进行缩放
		ratio1 = MaxWidth / img1Width // 图片1的缩放比例
		ratio2 = MaxWidth / img2Width // 图片2的缩放比例

		// 原宽度 * 比例 = 新宽度
		return int(img1Width * ratio1), int(img2Width * ratio2)
	}

	// 如果最小宽度小于600，那么需要将较大的图片缩放，使得两张图片的宽度一致
	if minWidth == img1Width {
		ratio2 = minWidth / img2Width // 图片2的缩放比例
		return img1W, int(img2Width * ratio2)
	}

	ratio1 = minWidth / img1Width // 图片1的缩放比例
	return int(img1Width * ratio1), img2W
}

func Merge(p1, p2, p3 string) {
	file1, _ := os.Open(p1) //打开图片1
	file2, _ := os.Open(p2) //打开图片2
	defer file1.Close()
	defer file2.Close()

	// image.Decode 图片
	var (
		img1, img2 image.Image
		err        error
	)
	if img1, _, err = image.Decode(file1); err != nil {
		return
	}
	if img2, _, err = image.Decode(file2); err != nil {
		return
	}
	b1 := img1.Bounds()
	b2 := img2.Bounds()
	new1W, new2W := fixSize(b1.Max.X, b2.Max.X)

	// 调用resize库进行图片缩放(高度填0，resize.Resize函数中会自动计算缩放图片的宽高比)
	m1 := resize.Resize(uint(new1W), 0, img1, resize.Lanczos3)
	m2 := resize.Resize(uint(new2W), 0, img2, resize.Lanczos3)

	// 将两个图片合成一张
	newWidth := m1.Bounds().Max.X                                                                          //新宽度 = 随意一张图片的宽度
	newHeight := m1.Bounds().Max.Y + m2.Bounds().Max.Y                                                     // 新图片的高度为两张图片高度的和
	newImg := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight))                                        //创建一个新RGBA图像
	draw.Draw(newImg, newImg.Bounds(), m1, m1.Bounds().Min, draw.Over)                                     //画上第一张缩放后的图片
	draw.Draw(newImg, newImg.Bounds(), m2, m2.Bounds().Min.Sub(image.Pt(0, m1.Bounds().Max.Y)), draw.Over) //画上第二张缩放后的图片（这里需要注意Y值的起始位置）

	// 保存文件
	imgfile, _ := os.Create(p3)
	defer imgfile.Close()
	jpeg.Encode(imgfile, newImg, &jpeg.Options{100})
}

func (that *Image) SaveToJpeg(quality int, savePath string) (err error) {
	//var img image.Image
	//var f bytes.Buffer
	f, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = jpeg.Encode(f, that.Image, &jpeg.Options{Quality: quality})
	if err != nil {
		return err
	}

	return nil
}

func (that *Image) SaveToWebP(quality float32, savePath string) (err error) {
	var img image.Image
	reader := bytes.NewReader(that.Data)
	lossLess := false //是否无损压缩
	Exact := false    //透明部分消失
	switch that.ImageType {
	case ImageJpeg:
		img, _ = jpeg.Decode(reader)
		break
	case ImagePng:
		img, _ = png.Decode(reader)
		lossLess = true
		Exact = true
		break
	case ImageBmp:
		img, _ = bmp.Decode(reader)
		break
	case ImageGif:
		//return that.gifToWebP(that.Data, quality)
	case ImageWebp:
		//return that.Data, nil
	}
	if img == nil {
		msg := "image file " + that.FilePath + " is corrupted or not supported"
		err = errors.New(msg)
		return err
	}

	if err = webp.Save(savePath, img, &webp.Options{Lossless: lossLess, Exact: Exact, Quality: quality}); err != nil {
		return err
	}

	return nil
}

// ToWebP /**
func (that *Image) ToWebP(quality float32) (out []byte, err error) {
	var img image.Image
	reader := bytes.NewReader(that.Data)
	lossLess := false //是否无损压缩
	Exact := false    //透明部分消失
	switch that.ImageType {
	case ImageJpeg:
		img, _ = jpeg.Decode(reader)
		break
	case ImagePng:
		img, _ = png.Decode(reader)
		lossLess = true
		Exact = true
		break
	case ImageBmp:
		img, _ = bmp.Decode(reader)
		break
	case ImageGif:
		return that.gifToWebP(that.Data, quality)
	case ImageWebp:
		return that.Data, nil
	}
	if img == nil {
		msg := "image file " + that.FilePath + " is corrupted or not supported"
		err = errors.New(msg)
		return nil, err
	}
	var buf bytes.Buffer
	if err = webp.Encode(&buf, img, &webp.Options{Lossless: lossLess, Exact: Exact, Quality: quality}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

/*
*
gif 转webP
*/
func (that *Image) gifToWebP(gifBin []byte, quality float32) (webPBin []byte, err error) {
	//todo
	return nil, nil
	//webpanim := giftowebp.NewWebpAnimation(111, 111, 1)
	//webpanim.WebPAnimEncoderOptions.SetKmin(9)
	//webpanim.WebPAnimEncoderOptions.SetKmax(17)
	//defer webpanim.ReleaseMemory() // dont forget call this or you will have memory leaks
	//webpConfig := giftowebp.NewWebpConfig()
	//webpConfig.SetLossless(1)
	//
	//timeline := 0
	//
	//for i, img := range gif.Image {
	//
	//	err = webpanim.AddFrame(img, timeline, webpConfig)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	timeline += gif.Delay[i] * 10
	//}
	//err = webpanim.AddFrame(nil, timeline, webpConfig)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = webpanim.Encode(&buf) // encode animation and write result bytes in buffer
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//ioutil.WriteFile("animation.webp", buf.Bytes(), 0777)
	//////////////
	//converter := giftowebp.NewConverter()
	//converter.LoopCompatibility = false
	////0 有损压缩  1无损压缩
	//converter.WebPConfig.SetLossless(0)
	////压缩速度  0-6  0最快 6质量最好
	//converter.WebPConfig.SetMethod(0)
	//converter.WebPConfig.SetQuality(quality)
	////搞不懂什么意思,例子是这样用的
	//converter.WebPAnimEncoderOptions.SetKmin(9)
	//converter.WebPAnimEncoderOptions.SetKmax(17)
	//
	//return converter.Convert(gifBin)
}

// MakeWebpThumbnail /**
func (that *Image) MakeWebpThumbnail(width int, height int, quality float32) (err error) {

	var img image.Image
	var desWidth int
	var desHeight int
	lossLess := false //是否无损压缩
	Exact := false    //透明部分消失
	reader := bytes.NewReader(that.Data)
	reader2 := bytes.NewReader(that.Data)
	switch that.ImageType {
	case ImageJpeg:
		img, _ = jpeg.Decode(reader)
		img2, _ := jpeg.DecodeConfig(reader2)
		desWidth = img2.Width
		desHeight = img2.Height
		break
	case ImagePng:
		img, _ = png.Decode(reader)
		img2, _ := png.DecodeConfig(reader2)
		desWidth = img2.Width
		desHeight = img2.Height
		lossLess = true
		Exact = true
		break
	case ImageBmp:
		img, _ = bmp.Decode(reader)
		img2, _ := bmp.DecodeConfig(reader2)
		desWidth = img2.Width
		desHeight = img2.Height
		break
	case ImageGif:
		//gifData, err := that.resizeGif(width, height)
		//if err != nil {
		//	return err
		//}
		//return that.gifToWebP(gifData, quality)
	case ImageWebp:
		img, _ = webp.Decode(reader)
		desWidth = that.Width
		desHeight = that.Height
		break
	}
	if img == nil {
		msg := "image file " + that.FilePath + " is corrupted or not supported"
		err = errors.New(msg)
		return err
	}
	w, h := that.calculateRatioFit(desWidth, desHeight, width, height)
	//var buf bytes.Buffer
	savePath := GetThumbPath(that.FilePath)
	fmt.Println(savePath)
	imgfile, _ := os.Create(savePath)
	defer imgfile.Close()
	m := resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	if err = webp.Encode(imgfile, m, &webp.Options{Lossless: lossLess, Exact: Exact, Quality: quality}); err != nil {
		return err
	}
	return nil
}

/*
*
改变gif的长宽
*/
func (that *Image) resizeGif(width int, height int) (out []byte, err error) {
	reader := bytes.NewReader(that.Data)
	im, err := gif.DecodeAll(reader)
	if err != nil {
		return nil, err
	}
	// reset the gif width and height
	im.Config.Width = width
	im.Config.Height = height

	firstFrame := im.Image[0].Bounds()
	img := image.NewRGBA(image.Rect(0, 0, firstFrame.Dx(), firstFrame.Dy()))

	// resize frame by frame
	for index, frame := range im.Image {
		b := frame.Bounds()
		draw.Draw(img, b, frame, b.Min, draw.Over)
		im.Image[index] = that.imageToPaletted(resize.Resize(uint(width), uint(height), img, resize.NearestNeighbor))
	}
	var buf bytes.Buffer
	err = gif.EncodeAll(&buf, im)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (that *Image) imageToPaletted(img image.Image) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, palette.Plan9)
	draw.FloydSteinberg.Draw(pm, b, img, image.ZP)
	return pm
}
