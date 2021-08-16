package assets

import (
	"bytes"
	"compress/gzip"
	"image"
	"io/ioutil"
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
