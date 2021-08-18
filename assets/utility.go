package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"image"
	"io/ioutil"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
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

func LoadImage(data []byte) (*ebiten.Image, error) {
	hash := getHash(data)

	lock.Lock()
	defer lock.Unlock()
	eImg, ok := imageCache[hash]
	if ok {
		return eImg, nil
	}

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

	eImg = ebiten.NewImageFromImage(img)

	imageCache[hash] = eImg

	return eImg, nil
}
