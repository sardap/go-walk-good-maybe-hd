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
	"github.com/nfnt/resize"

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
	scale := int(t.FieldByName("ScaleMultiplier").Int())

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
	img = resize.Resize(uint(img.Bounds().Dx()*scale), uint(img.Bounds().Dy()*scale), img, resize.NearestNeighbor)

	eImg = ebiten.NewImageFromImage(img)

	imageCache[hash] = eImg

	return eImg, nil
}

func LoadSound(asset interface{}) (data []byte, sampleRate int, soundType SoundType) {
	t := reflect.ValueOf(asset)

	sampleRate = int(t.FieldByName("SampleRate").Int())
	data = []byte(t.FieldByName("Data").String())
	soundType = SoundType(t.FieldByName("SoundType").Int())

	return
}
