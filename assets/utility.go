package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"image"
	"io/ioutil"
	"reflect"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"

	_ "image/png"
)

var (
	imageCache map[[16]byte]*ebiten.Image
	lock       *sync.Mutex
)

func init() {
	imageCache = make(map[[16]byte]*ebiten.Image)
	lock = &sync.Mutex{}
}

func getHash(data []byte) [16]byte {
	return md5.Sum(data)
}

func LoadEbitenImage(asset interface{}) (*ebiten.Image, error) {
	t := reflect.ValueOf(asset)

	compressed := t.FieldByName("Compressed").Bool()
	data := []byte(t.FieldByName("Data").String())

	hash := getHash(data)

	lock.Lock()
	defer lock.Unlock()
	eImg, ok := imageCache[hash]
	if ok {
		return eImg, nil
	}

	if compressed {
		zr, _ := gzip.NewReader(bytes.NewReader(data))
		defer zr.Close()
		var err error
		data, err = ioutil.ReadAll(zr)
		if err != nil {
			panic(err)
		}
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	eImg = ebiten.NewImageFromImage(img)

	imageCache[hash] = eImg

	return eImg, nil
}
