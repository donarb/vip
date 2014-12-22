package fetch

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
)

func needsRotation(src io.Reader) int {
	metadata, err := exif.Decode(src)
	if err != nil {
		return 0
	}

	orientation, err := metadata.Get(exif.Orientation)
	if err != nil {
		return 0
	}

	switch orientation.String() {
	case "6":
		return 270
	case "3":
		return 180
	case "8":
		return 90
	default:
		return 0
	}

}

func GetRotatedImage(src io.Reader) (image.Image, string, error) {
	raw, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, "", err
	}

	data := bytes.NewReader(raw)

	image, format, err := image.Decode(data)
	if err != nil {
		return nil, "", err
	}

	if _, err := data.Seek(0, 0); err != nil {
		return nil, "", err
	}

	angle := needsRotation(data)
	switch angle {
	case 90:
		image = imaging.Rotate90(image)
	case 180:
		image = imaging.Rotate180(image)
	case 270:
		image = imaging.Rotate270(image)
	}

	return image, format, nil
}

func Resize(src io.Reader, c *CacheContext) (io.Reader, error) {
	image, format, err := image.Decode(src)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	factor := float64(c.Width) / float64(image.Bounds().Size().X)
	height := int(float64(image.Bounds().Size().Y) * factor)

	image = imaging.Resize(image, c.Width, height, imaging.Linear)

	buf := new(bytes.Buffer)

	switch format {
	case "jpeg":
		err = jpeg.Encode(buf, image, nil)
	case "png":
		err = png.Encode(buf, image)
	}

	return buf, err
}

func CenterCrop(src io.Reader, c *CacheContext) (io.Reader, error) {
	image, format, err := image.Decode(src)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	height := image.Bounds().Size().Y
	width := image.Bounds().Size().X

	if width < height {
		image = imaging.CropCenter(image, width, width)
	} else if width > height {
		image = imaging.CropCenter(image, height, height)
	} else {
		image = imaging.CropCenter(image, width, height)
	}

	buf := new(bytes.Buffer)

	switch format {
	case "jpeg":
		err = jpeg.Encode(buf, image, nil)
	case "png":
		err = png.Encode(buf, image)
	}

	return buf, err
}
