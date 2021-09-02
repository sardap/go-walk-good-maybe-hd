package assets

import (
	"bytes"
	"compress/gzip"
	"testing"
	"time"

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
	t.Parallel()

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

	img, err := LoadEbitenImage(asset)
	assert.NoError(t, err)
	assert.Equal(t, int(16), img.Bounds().Max.X)

	startTime := time.Now()
	for i := 0; i < 1000; i++ {
		delete(imageCache, getHash([]byte(asset.Data)))
		LoadEbitenImage(asset)
	}
	delta := time.Since(startTime)

	startTime = time.Now()
	for i := 0; i < 1000; i++ {
		LoadEbitenImage(asset)
	}
	assert.Less(t, time.Since(startTime), delta/15)
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

	img, err = LoadEbitenImage(asset)
	assert.NoError(t, err)
	assert.Equal(t, int(16), img.Bounds().Max.X)
}

func TestLoadSound(t *testing.T) {
	t.Parallel()

	asset := struct {
		SampleRate int
		Data       string
		SoundType  SoundType
	}{
		SampleRate: 10312312,
		Data:       "looky here",
		SoundType:  SoundTypeWav,
	}

	data, sr, soundType := LoadSound(asset)
	assert.Equal(t, int(10312312), sr)
	assert.Equal(t, []byte("looky here"), data)
	assert.Equal(t, SoundTypeWav, soundType)
}
