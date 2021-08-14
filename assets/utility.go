package assets

import (
	"bytes"
	"compress/gzip"
	"image"
	"io/ioutil"

	"github.com/nfnt/resize"
)

func LoadImage(data []byte) (image.Image, error) {

	zr, _ := gzip.NewReader(bytes.NewReader(data))
	defer zr.Close()
	uncompressed, err := ioutil.ReadAll(zr)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewReader(uncompressed))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ScaleImage(img image.Image) image.Image {
	return resize.Resize(uint(img.Bounds().Dx())*8, uint(img.Bounds().Dy())*8, img, resize.NearestNeighbor)
}
