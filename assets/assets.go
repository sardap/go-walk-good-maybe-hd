package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
	"github.com/sardap/walk-good-maybe-hd/common"
)

type SoundType int

const (
	SoundTypeMp3 SoundType = 0
	SoundTypeWav SoundType = 1
)

var (
	imageCache map[[16]byte]*ebiten.Image
	lock       *sync.Mutex
	assetsPath string = "out"
)

func SetAssetsPath(path string) {
	assetsPath = path
}

func ClearImageCache() {
	imageCache = make(map[[16]byte]*ebiten.Image)
}

func DeleteImageCache(hash [16]byte) {
	delete(imageCache, hash)
}

func init() {
	imageCache = make(map[[16]byte]*ebiten.Image)
	lock = &sync.Mutex{}
}

func getAsset(path string) (f io.ReadCloser) {
	f, err := os.Open(filepath.Join(assetsPath, path))
	if err != nil {
		resp, err := http.Get(Remote + filepath.ToSlash(path))
		if err != nil {
			panic(err)
		}

		f = resp.Body
	}

	return
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

func LoadEbitenImageAsset(asset interface{}) (*ebiten.Image, error) {
	return LoadEbitenImageColorSwap(asset, nil)
}

func LoadEbitenImageRaw(imageData []byte) (*ebiten.Image, error) {
	hash := getImageHash(imageData, nil)

	lock.Lock()
	defer lock.Unlock()
	eImg, ok := imageCache[hash]
	if ok {
		return eImg, nil
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	eImg = ebiten.NewImageFromImage(img)
	imageCache[hash] = eImg

	return eImg, nil
}

func LoadSound(asset interface{}) (data []byte, sampleRate int, soundType SoundType) {
	t := reflect.ValueOf(asset)

	if field := t.FieldByName("Path"); field.IsValid() {
		data, _ = ioutil.ReadAll(getAsset(field.String()))
	} else {
		data = []byte(t.FieldByName("Data").String())
	}

	sampleRate = int(t.FieldByName("SampleRate").Int())
	soundType = SoundType(t.FieldByName("SoundType").Int())

	return
}

func LoadKaraokeSession(path string) (session *common.KaraokeSession) {

	f := getAsset(filepath.Join("karaoke", path))
	defer f.Close()

	zr, _ := gzip.NewReader(f)
	defer zr.Close()
	data, err := ioutil.ReadAll(zr)
	if err != nil {
		panic(err)
	}

	session = &common.KaraokeSession{}
	toml.Decode(string(data), session)

	return
}

func LoadKaraokeIndex() *common.KaraokeIndex {
	r := getAsset(filepath.Join("karaoke", "index.toml"))
	defer r.Close()

	karaokeIndex := &common.KaraokeIndex{}
	dec := toml.NewDecoder(r)
	dec.Decode(karaokeIndex)

	return karaokeIndex
}
