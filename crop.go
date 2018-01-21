package main

import (
	"bytes"
	"image"
	"image/color"

	_ "image/gif"
	_ "image/jpeg"
	"image/png"
)

const (
	SIZE = 480
)

type CroppedPhoto struct {
	Original image.Image
}

func (cp CroppedPhoto) ColorModel() color.Model {
	return cp.Original.ColorModel()
}

func (cp CroppedPhoto) Bounds() image.Rectangle {
	return image.Rect(0, 0, SIZE, SIZE)
}

func (cp CroppedPhoto) At(x, y int) color.Color {
	// Get original bounds
	oRect := cp.Original.Bounds()

	// Get new top left corner
	nx := (oRect.Max.X-oRect.Min.X)/2 - (SIZE / 2) + x
	ny := (oRect.Max.Y-oRect.Min.Y)/2 - (SIZE / 2) + y

	// If it's part of the original image just grab that
	if inBounds(nx, ny, oRect) {
		return cp.Original.At(nx, ny)
	} else {
		// Otherwise just white :shrug:
		return color.White
	}
}

func inBounds(x, y int, box image.Rectangle) bool {
	return x >= box.Min.X &&
		x <= box.Max.X &&
		y >= box.Min.Y &&
		y <= box.Max.Y
}

func cropBytesToBytes(origB []byte) ([]byte, error) {
	readBuf := bytes.NewReader(origB)
	origImg, _, err := image.Decode(readBuf)
	if err != nil {
		return nil, err
	}

	croppedImg := CroppedPhoto{origImg}
	var writeBuf bytes.Buffer
	err = png.Encode(&writeBuf, croppedImg)
	if err != nil {
		return nil, err
	}

	return writeBuf.Bytes(), nil
}
