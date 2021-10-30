package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"

	"github.com/h2non/bimg"
	"github.com/nfnt/resize"
	"golang.org/x/image/draw"
)

const (
	imageWidth   = 640
	imageHeight  = 360
	imageQuality = 90
)

func ResizeImage(file io.Reader) (*bytes.Buffer, error) {
	return bimgResizeImage(file)
}

//github.com/h2non/bimg
func bimgResizeImage(file io.Reader) (*bytes.Buffer, error) {
	options := bimg.Options{
		Width:        imageWidth,
		Height:       imageHeight,
		Crop:         false,
		Extend:       bimg.ExtendWhite,
		Interpolator: bimg.Bilinear,
		Gravity:      bimg.GravityCentre,
		Quality:      imageQuality,
	}
	inBuf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	buf, err := bimg.Resize(inBuf, options)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buf), err
}

//"golang.org/x/image/draw"
func drawResizeImage(file io.Reader) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	img, err := jpeg.Decode(file)
	if err != nil {
		return buf, err
	}
	newImg := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.NearestNeighbor.Scale(newImg, newImg.Rect, img, img.Bounds(), draw.Over, nil)

	jpeg.Encode(buf, newImg, nil)
	return buf, nil
}

//"github.com/nfnt/resize"
func nfntResizeImage(file io.Reader) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	img, err := jpeg.Decode(file)
	if err != nil {
		return buf, err
	}
	newImg := resize.Resize(imageWidth, imageHeight, img, resize.Lanczos3)
	jpeg.Encode(buf, newImg, nil)
	return buf, nil
}
