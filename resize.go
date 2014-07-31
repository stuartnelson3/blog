package main

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func CreateImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	path := "./public/img/" + header.Filename
	img, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer img.Close()

	var decodedImage image.Image
	switch filepath.Ext(header.Filename) {
	case ".jpg", ".jpeg":
		decodedImage, err = jpeg.Decode(file)
	case ".png":
		decodedImage, err = png.Decode(file)
	case ".gif":
		// Resizing gifs isn't working, so just copy the file.
		io.Copy(img, file)
		return path, nil
	default:
		return "", errors.New("Bad filetype")
	}

	if err != nil {
		return "", err
	}

	if b := decodedImage.Bounds(); b.Dx() > 600 {
		scale := float64(600) / float64(b.Dx())
		width, height := Scale(b, scale)

		decodedImage = resize.Resize(uint(width), uint(height), decodedImage, resize.Lanczos3)
	}

	jpeg.Encode(img, decodedImage, nil)
	return path, nil
}

func Scale(image image.Rectangle, scale float64) (width, height uint) {
	return uint(float64(image.Dx()) * scale), uint(float64(image.Dy()) * scale)
}
