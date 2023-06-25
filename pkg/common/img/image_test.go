package app

import (
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"image/jpeg"
	"os"
	"strings"
	"testing"
)

const jpegRoot = "/Users/wan/Pictures/blogImage/jpeg/"
const webpRoot = "/Users/wan/Pictures/blogImage/webp/"
const webps = "/Users/wan/Pictures/blogImage/webp/%d.webp"
const jpegs = "/Users/wan/Pictures/blogImage/jpeg/%d.jpeg"
const jpegs_thumb = "/Users/wan/Pictures/blogImage/jpeg/thumb/%d.jpeg"

func HandleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestJpegToWebpThumbnail(t *testing.T) {
	path := "/Users/wan/Pictures/blogImage/poetize.jpg"
	image := Image{}
	_ = image.Open(path)

	savePath := "/Users/wan/Pictures/blogImage/poetize.webp"
	_ = image.SaveToWebP(75, savePath)

	image2 := Image{}
	_ = image2.Open(savePath)

	_ = image2.MakeWebpThumbnail(120, 120, 80)
}

func TestIndex(t *testing.T) {

}

func TestWebpToThumbnail(t *testing.T) {
	var webpPath string
	files, _ := os.ReadDir(webpRoot)

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".webp") {
			continue
		}

		webpPath = webpRoot + f.Name()
		if _, err := os.Stat(webpPath); err != nil {
			continue
		}
		if f.IsDir() {
			continue
		}

		image := Image{}
		if image.Open(webpPath) != nil {
			fmt.Println(webpPath)
		}
		HandleError(image.MakeWebpThumbnail(120, 120, 80))
	}

}

func TestJpegToThumbnail(t *testing.T) {
	var jpegPath string
	files, _ := os.ReadDir(jpegRoot)

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".jpeg") {
			continue
		}

		jpegPath = jpegRoot + f.Name()
		if _, err := os.Stat(jpegPath); err != nil {
			continue
		}
		if f.IsDir() {
			continue
		}

		MakeJpegThumbnail(jpegPath)
	}

}

func TestMerge(t *testing.T) {
	Merge(jpegRoot+"thumb/1.jpeg", jpegRoot+"thumb/2.jpeg", jpegRoot+"thumb/3.jpeg")
}

func TestFileCount(t *testing.T) {
	files, _ := os.ReadDir("/Users/wan/Pictures/download/jpeg/")
	for _, f := range files {
		info, _ := f.Info()
		fmt.Println(f.Name(), info.Size(), f.Type())
	}
	fmt.Println(len(files))
}

func TestJpegToThumb(t *testing.T) {
	var jpegPath string
	files, _ := os.ReadDir(jpegRoot)
	//for _, f := range files {
	//	fmt.Println(f.Name())
	//}

	for _, f := range files {
		//webpPath = fmt.Sprintf(webps, i)
		if !strings.HasSuffix(f.Name(), ".jpeg") {
			continue
		}

		jpegPath = jpegRoot + f.Name()
		if _, err := os.Stat(jpegPath); err != nil {
			continue
		}
		if f.IsDir() {
			continue
		}

		image := Image{}
		if err := image.Open(jpegPath); err != nil {
			fmt.Println(f.Name(), err.Error())
		}

		width, height := 480, 480
		if err := image.JpegToThumb(width, height, 75, jpegRoot+"thumb/"+f.Name()); err != nil {
			fmt.Println(f.Name(), err.Error())
		}

		image.Reset()
	}
}

func TestImage_Convert(t *testing.T) {
	//var filePaths []string
	var webpPath, jpegPath string
	for i := 1; i < 37; i++ {
		//filePaths = append(filePaths, fmt.Sprintf(basePath, i))
		webpPath = fmt.Sprintf(webps, i)
		image := Image{}
		if err := image.Open(webpPath); err != nil {
			fmt.Println("\n", i)
			//panic(err.Error())
			continue
		}

		jpegPath = fmt.Sprintf(jpegs, i)
		if err := image.SaveToJpeg(50, jpegPath); err != nil {
			fmt.Println("\n", i)
			//panic(err.Error())
			continue
		}

		image.Reset()
	}

}

func TestImage_Open(t *testing.T) {

	img := Image{}
	filePath := "/Users/wan/Pictures/download/4_jpg.jpg"

	err := img.Open(filePath)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(img.Data)
	img2, _ := jpeg.Decode(reader)

	//bs, err := img.ToWebP(1)
	//if err != nil {
	//	panic(err)
	//}

	err = webp.Save("/Users/wan/Pictures/download/4_jpg.jpg.webp", img2, &webp.Options{Quality: 50})
	if err != nil {
		panic(err)
	}
}
