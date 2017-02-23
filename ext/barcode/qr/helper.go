package qr

import (
	"bytes"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/henrylee2cn/faygo/ext/barcode"
)

// 常见二维码边长
const (
	L200  int = 200
	L300  int = 300
	L500  int = 500
	L600  int = 600
	L800  int = 800
	L1000 int = 1000
	L1500 int = 1500
	L2000 int = 2000
)

// 输出二维码对象
// ecl为容错性
func Img(content string, ecl ErrorCorrectionLevel, sideLength int) (barcode.Barcode, error) {
	img, err := Encode(content, ecl, Auto)
	if err != nil {
		return img, err
	}

	img, err = barcode.Scale(img, sideLength, sideLength)
	return img, err
}

// 输出png格式图片字节流
func Png(content string, ecl ErrorCorrectionLevel, sideLength int) (bytes.Buffer, error) {
	var data bytes.Buffer
	img, err := Img(content, ecl, sideLength)
	if err != nil {
		return data, err
	}

	err = png.Encode(&data, img)
	return data, err
}

// 输出gif格式图片字节流
func Gif(content string, ecl ErrorCorrectionLevel, sideLength int, o *gif.Options) (bytes.Buffer, error) {
	var data bytes.Buffer
	img, err := Img(content, ecl, sideLength)
	if err != nil {
		return data, err
	}
	err = gif.Encode(&data, img, o)
	return data, err
}

// 输出jpeg格式图片字节流
func Jpeg(content string, ecl ErrorCorrectionLevel, sideLength int, o *jpeg.Options) (bytes.Buffer, error) {
	var data bytes.Buffer
	img, err := Img(content, ecl, sideLength)
	if err != nil {
		return data, err
	}
	err = jpeg.Encode(&data, img, o)
	return data, err
}
