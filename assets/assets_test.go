package assets_test

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"errors"
	"image"
	"image/color"
	"os"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/stretchr/testify/assert"
)

const (
	testImg = "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x10\x00\x00\x00\b\b\x02\x00\x00\x00\u007f\x14\xe8\xc0\x00\x00\x00\u007fIDATx\x9cb\xa9\xad\xcde\xc0\x01\x94\x85\xf8\xd0D\xee\xbe\xfbĂK5\x03\x03\x83\x94$\x13\x86\x06\x06\x90\x06\xf6\x1f_\xd1$~rpC\x18\x95\x1f\x9b\x18\x18\x18\f\xf5A\xec\xb0\au\f\f\f\xe8f`\x05\xe7/\x82\x10\x04@5\xc8HI\xa21p\x01\x16\x88\xa2'ϞC\x94B\x18w\xdf}§\x01\xa2\xe8ɳ\xe7p\xcd\f0?\xb4\xf3\xd7a\xb7\x01\xcdU\x10\x1b\x9e=\xff\x87i\x03 \x00\x00\xff\xff\x80O*\xbb\x1c7\xeb\xc5\x00\x00\x00\x00IEND\xaeB`\x82"
)

func compress(data []byte) []byte {
	compressed := &bytes.Buffer{}
	func() {
		zw := gzip.NewWriter(compressed)
		defer zw.Close()

		_, err := zw.Write(data)
		if err != nil {
			panic(err)
		}
	}()

	return compressed.Bytes()
}

func TestLoadEbitenImage(t *testing.T) {

	var asset struct {
		Compressed      bool
		Data            string
		ScaleMultiplier int
	}

	// Uncompressed
	asset = struct {
		Compressed      bool
		Data            string
		ScaleMultiplier int
	}{
		Compressed:      false,
		Data:            testImg,
		ScaleMultiplier: 1,
	}

	img, err := assets.LoadEbitenImageAsset(asset)
	assert.NoError(t, err)
	assert.Equal(t, int(16), img.Bounds().Max.X)

	startTime := time.Now()
	for i := 0; i < 1000; i++ {
		assets.DeleteImageCache(md5.Sum([]byte(asset.Data)))
		assets.LoadEbitenImageAsset(asset)
	}
	delta := time.Since(startTime)

	startTime = time.Now()
	for i := 0; i < 1000; i++ {
		assets.LoadEbitenImageAsset(asset)
	}
	assert.Less(t, time.Since(startTime), delta/5)
	assert.NoError(t, err)

	// Compressed
	asset = struct {
		Compressed      bool
		Data            string
		ScaleMultiplier int
	}{
		Compressed:      true,
		ScaleMultiplier: 1,
		Data:            string(compress([]byte(testImg))),
	}

	img, err = assets.LoadEbitenImageAsset(asset)
	assert.NoError(t, err)
	assert.Equal(t, int(16), img.Bounds().Max.X)
}

func colorToRGBA(r, g, b, a uint32) color.RGBA {
	return color.RGBA{R: byte(r), G: byte(g), B: byte(b), A: byte(a)}
}

func TestLoadEbitenImageColorSwap(t *testing.T) {

	asset := struct {
		Compressed      bool
		Data            string
		ScaleMultiplier int
	}{
		Compressed:      false,
		Data:            testImg,
		ScaleMultiplier: 1,
	}

	originalImg, _, _ := image.Decode(bytes.NewBuffer([]byte(testImg)))
	colorMap := map[color.RGBA]color.RGBA{}
	originalColor := colorToRGBA(originalImg.At(0, 0).RGBA())
	newColor := color.RGBA{R: 255, A: 255}
	colorMap[originalColor] = newColor

	updates := []struct {
		X int
		Y int
	}{}

	for y := originalImg.Bounds().Min.Y; y < originalImg.Bounds().Max.Y; y++ {
		for x := originalImg.Bounds().Min.X; x < originalImg.Bounds().Max.X; x++ {
			if originalImg.At(x, y) == originalColor {
				updates = append(updates, struct {
					X int
					Y int
				}{x, y})
			}
		}
	}

	img, err := assets.LoadEbitenImageColorSwap(asset, colorMap)
	assert.NoError(t, err)

	for _, update := range updates {
		updated := img.At(update.X, update.Y)
		assert.Equalf(t, updated, newColor, "missmatch at X:%d, Y:%d", update.X, update.Y)
	}
}

func TestLoadSound(t *testing.T) {
	t.Parallel()

	asset := struct {
		SampleRate int
		Data       string
		SoundType  assets.SoundType
	}{
		SampleRate: 10312312,
		Data:       "looky here",
		SoundType:  assets.SoundTypeWav,
	}

	data, sr, soundType := assets.LoadSound(asset)
	assert.Equal(t, int(10312312), sr)
	assert.Equal(t, []byte("looky here"), data)
	assert.Equal(t, assets.SoundTypeWav, soundType)
}

func TestLoadKaraoke(t *testing.T) {
	t.Parallel()

	data := assets.LoadKaraoke(assets.KaraokePdRock01)
	assert.NotNil(t, data)
}

type testGame struct {
	m    *testing.M
	code int
}

var (
	errRegularTermination = errors.New("regular termination")
)

func (g *testGame) Update() error {
	g.code = g.m.Run()
	return errRegularTermination
}

func (*testGame) Draw(screen *ebiten.Image) {
}

func (*testGame) Layout(int, int) (int, int) {
	return 300, 300
}

func TestMain(m *testing.M) {
	g := &testGame{
		m: m,
	}
	if err := ebiten.RunGame(g); err != nil && err != errRegularTermination {
		panic(err)
	}
	os.Exit(g.code)
}
