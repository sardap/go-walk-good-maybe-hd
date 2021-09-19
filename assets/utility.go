package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"image"
	"io/ioutil"
	"reflect"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"

	"image/color"
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

func getImageHash(data []byte, clrMap map[color.RGBA]color.RGBA) [16]byte {
	result := md5.Sum(data)
	if clrMap != nil {
		// LOOK: slowpoint
		buf := &bytes.Buffer{}
		je := json.NewEncoder(buf)

		for key, value := range clrMap {
			je.Encode(key)
			je.Encode(value)
		}

		buf.Write(result[:])
		result = md5.Sum(buf.Bytes())
 	}

	return result
}

func LoadEbitenImageColorSwap(asset interface{}, clrMap map[color.RGBA]color.RGBA) (*ebiten.Image, error) {
	t := reflect.ValueOf(asset)

	data := []byte(t.FieldByName("Data").String())

	hash := getImageHash(data, clrMap)

	lock.Lock()
	defer lock.Unlock()
	eImg, ok := imageCache[hash]
	if ok {
		return eImg, nil
	}

	compressed := t.FieldByName("Compressed").Bool()
	scale := int(t.FieldByName("ScaleMultiplier").Int())

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
	if clrMap != nil {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				mappedClr, ok := clrMap[eImg.At(x, y).(color.RGBA)]
				if ok {
					eImg.Set(x, y, mappedClr)
				}
			}
		}
	}

	imageCache[hash] = eImg

	return eImg, nil
}

func LoadEbitenImage(asset interface{}) (*ebiten.Image, error) {
	return LoadEbitenImageColorSwap(asset, nil)
}

func LoadSound(asset interface{}) (data []byte, sampleRate int, soundType SoundType) {
	t := reflect.ValueOf(asset)

	sampleRate = int(t.FieldByName("SampleRate").Int())
	data = []byte(t.FieldByName("Data").String())
	soundType = SoundType(t.FieldByName("SoundType").Int())

	return
}
