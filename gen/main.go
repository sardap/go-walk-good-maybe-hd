package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/mjibson/go-dsp/wav"
	"github.com/tcolgate/mp3"

	_ "github.com/oov/psd"
)

const warning = "AUTO GENERATED CODE DO NOT EDIT REFER TO gen/main.go"

var workspacePath string

func ptr(typeName string) string {
	return fmt.Sprintf("*%s", typeName)
}

func genComInterface(f *jen.File, name string) {
	rawName := strings.Replace(name, "Component", "", 1)

	f.Type().Id(fmt.Sprintf("%sFace", rawName)).Interface(
		jen.Id(fmt.Sprintf("Get%s", name)).Params().Qual("", ptr(name)),
	)
}

func genComGetFunction(f *jen.File, name string) {
	firstLetter := strings.ToLower(string(name[0]))

	f.Func().Params(
		jen.Id(firstLetter).Id(ptr(name)),
	).Id(fmt.Sprintf("Get%s", name)).Params().Qual("", ptr(name)).Block(
		jen.Return(jen.Id(firstLetter)),
	)
}

func genComCode() {
	fmt.Printf("Generating components\n")

	jf := jen.NewFile("components")

	jf.Comment(warning)
	jf.Line()

	err := filepath.Walk(filepath.Join(workspacePath, "components"),
		func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), "Component.go") {
				return err
			}

			name := info.Name()
			name = strings.TrimSuffix(name, ".go")

			genComGetFunction(jf, name)
			jf.Line()
			genComInterface(jf, name)
			jf.Line()

			return nil
		})
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(workspacePath, "components", "gen.go"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	jf.Render(f)
}

func compressAsset(data []byte) ([]byte, int) {
	compressed := &bytes.Buffer{}
	func() {
		zw := gzip.NewWriter(compressed)
		defer zw.Close()

		_, err := zw.Write(data)
		if err != nil {
			panic(err)
		}
	}()

	return compressed.Bytes(), len(data) - compressed.Len()
}

func genByteArray(jf *jen.File, data []byte, name string) {
	jf.Var().Id(name).Op("=").Index().Byte().Params(jen.Lit(string(data)))
	jf.Line()
}

type GraphicsOutput struct {
	Name            string
	Type            string
	Parts           []string
	FrameWidth      int
	ScaleMultiplier int
	File            string
}

func (g *GraphicsOutput) genImageAsset(jf *jen.File, img image.Image) {
	var fields []jen.Code
	fields = append(fields, jen.Id("Data").String())
	fields = append(fields, jen.Id("ScaleMultiplier").Int())
	fields = append(fields, jen.Id("Compressed").Bool())

	if g.FrameWidth > 0 {
		fields = append(fields, jen.Id("FrameWidth").Int())
	}

	buffer := &bytes.Buffer{}
	png.Encode(buffer, img)
	imgBytes := buffer.Bytes()

	// Copy this
	compressed, diff := compressAsset(imgBytes)
	useCompressed := diff > 40
	if useCompressed {
		fmt.Printf("Reduced %s by %d bytes\n", g.Name, diff)
		imgBytes = compressed
	} else {
		fmt.Printf("Not compressing %s too small\n", g.Name)
	}

	name := g.Name
	name = strcase.ToCamel(name)
	name = "Image" + name

	if g.Type == "tileSet" && len(g.Parts) == 0 {
		for i := 0; i < img.Bounds().Dx()/g.FrameWidth; i++ {
			g.Parts = append(g.Parts, fmt.Sprintf("Frame%d", i))
		}
	}

	if len(g.Parts) > 0 {
		name += "TileSet"
	}

	for i, part := range g.Parts {
		id := strcase.ToCamel("Index " + g.Name + " " + part)
		jf.Const().Id(id).Op("=").Lit(i)
	}

	if g.ScaleMultiplier == 0 {
		g.ScaleMultiplier = 1
	}

	jf.Var().Id(name).Op("=").Struct(fields...).BlockFunc(func(jf *jen.Group) {
		jf.Id("Data").Op(":").Lit(string(imgBytes)).Op(",")
		jf.Id("ScaleMultiplier").Op(":").Lit(int(g.ScaleMultiplier)).Op(",")

		if g.FrameWidth > 0 {
			jf.Id("FrameWidth").Op(":").Lit(g.FrameWidth * g.ScaleMultiplier).Op(",")
		}

		jf.Id("Compressed").Op(":").Lit(useCompressed).Op(",")
	})

	jf.Line()
}

func (g *GraphicsOutput) genImageAssetFromFile(jf *jen.File, path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	g.genImageAsset(jf, img)
}

func genImagesAssets(jf *jen.File, images []GraphicsOutput, assetsPath string) {
	for _, target := range images {
		target.genImageAssetFromFile(jf, filepath.Join(assetsPath, target.File))
	}
}

type MusicOutput struct {
	Name string
	File string
}

func getMp3SampleRate(path string) int {
	mp3F, _ := os.Open(path)
	defer mp3F.Close()
	d := mp3.NewDecoder(mp3F)
	var f mp3.Frame
	skipped := 0
	if err := d.Decode(&f, &skipped); err != nil {
		panic(err)
	}

	sampleRate := f.Header().SampleRate()
	if sampleRate == mp3.ErrInvalidSampleRate {
		panic("Invalid sample rate")
	}

	return int(sampleRate)

}

func (s *MusicOutput) genMusicAssetFromFile(jf *jen.File, path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	name = strcase.ToCamel(name)
	name = "Music" + name

	sampleRate := getMp3SampleRate(path)

	fields := []jen.Code{
		jen.Id("Data").String(),
		jen.Id("SampleRate").Int(),
		jen.Id("SoundType").Op("SoundType"),
	}

	jf.Var().Id(name).Op("=").Struct(fields...).BlockFunc(func(jf *jen.Group) {
		jf.Id("Data").Op(":").Lit(string(data)).Op(",")
		jf.Id("SampleRate").Op(":").Lit(sampleRate).Op(",")
		jf.Id("SoundType").Op(":").Op("SoundTypeMp3").Op(",")
	})
}

func genMusicAssets(jf *jen.File, music []MusicOutput, assetsPath string) {
	fmt.Printf("\nGenerating Music\n")

	for _, target := range music {
		fmt.Printf("...%s\n", target.Name)
		target.genMusicAssetFromFile(jf, filepath.Join(assetsPath, target.File))
	}
}

type SoundOutput struct {
	Name string
	File string
}

func getWavSampleRate(path string) int {
	wavF, _ := os.Open(path)
	defer wavF.Close()
	w, _ := wav.New(wavF)
	return int(w.SampleRate)
}

func (s *SoundOutput) genSoundAssetFromFile(jf *jen.File, path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	name = strcase.ToCamel(name)
	name = "Sound" + name

	sampleRate := getWavSampleRate(path)

	fields := []jen.Code{
		jen.Id("Data").String(),
		jen.Id("SampleRate").Int(),
		jen.Id("SoundType").Op("SoundType"),
	}

	jf.Var().Id(name).Op("=").Struct(fields...).BlockFunc(func(jf *jen.Group) {
		jf.Id("Data").Op(":").Lit(string(data)).Op(",")
		jf.Id("SampleRate").Op(":").Lit(sampleRate).Op(",")
		jf.Id("SoundType").Op(":").Op("SoundTypeWav").Op(",")
	})
}

func genSoundAssets(jf *jen.File, sounds []SoundOutput, assetsPath string) {
	fmt.Printf("\nGenerating Sounds\n")

	for _, target := range sounds {
		fmt.Printf("...%s\n", target.Name)
		target.genSoundAssetFromFile(jf, filepath.Join(assetsPath, target.File))
	}
}

type KaraokeOutput struct {
	Name string
	File string
}

type KaraokeInput struct {
	StartTime   int    `json:"start_time"`
	Duration    int    `json:"duration"`
	StartOffset int    `json:"start_off,omitempty"`
	Sound       string `json:"sound"`
}

type KaraokeBackground struct {
	Duration int    `json:"duration"`
	FadeIn   int    `json:"fade_in"`
	Image    string `json:"image"`
}

type KaraokeIn struct {
	Inputs      []KaraokeInput      `json:"inputs"`
	Backgrounds []KaraokeBackground `json:"backgrounds"`
	Music       string              `json:"music"`
	SampleRate  int                 `json:"sampleRate"`
}

func (s *KaraokeOutput) genKaraokeAssetFromFile(jf *jen.File, path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	karaoke := KaraokeIn{}
	err = json.Unmarshal(data, &karaoke)
	if err != nil {
		panic(err)
	}

	karaokePath := strings.TrimSuffix(path, filepath.Base(path))

	for i, background := range karaoke.Backgrounds {
		karaoke.Backgrounds[i].Image = func() string {
			f, err := os.Open(filepath.Join(karaokePath, background.Image))
			if err != nil {
				panic(err)
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				panic(err)
			}

			buffer := &bytes.Buffer{}
			if err := jpeg.Encode(buffer, img, &jpeg.Options{Quality: 30}); err != nil {
				panic(err)
			}

			return base64.StdEncoding.EncodeToString(buffer.Bytes())
		}()
	}

	timeElapsed := 0
	for i := 1; i < len(karaoke.Inputs); i++ {
		lastInput := &karaoke.Inputs[i-1]
		input := &karaoke.Inputs[i]
		input.StartTime = timeElapsed + lastInput.Duration + input.StartOffset
		input.StartOffset = 0
		timeElapsed += input.Duration
	}

	karaoke.SampleRate = getMp3SampleRate(filepath.Join(karaokePath, karaoke.Music))
	rawMusic, _ := ioutil.ReadFile(filepath.Join(karaokePath, karaoke.Music))
	karaoke.Music = string(base64.StdEncoding.EncodeToString(rawMusic))

	data, _ = json.Marshal(karaoke)

	compactJson := &bytes.Buffer{}
	json.Compact(compactJson, data)
	compressedJson, _ := compressAsset(compactJson.Bytes())

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	name = strcase.ToCamel(name)
	name = "Karaoke" + name

	fields := []jen.Code{
		jen.Id("JsonStr").String(),
	}

	jf.Var().Id(name).Op("=").Struct(fields...).BlockFunc(func(jf *jen.Group) {
		jf.Id("JsonStr").Op(":").Lit(string(compressedJson)).Op(",")
	})
}

func genKarokeAssets(jf *jen.File, karaoke []KaraokeOutput, assetsPath string) {
	fmt.Printf("\nGenerating Karaoke\n")

	for _, target := range karaoke {
		fmt.Printf("...%s\n", target.Name)
		target.genKaraokeAssetFromFile(jf, filepath.Join(assetsPath, target.File))
	}
}

func genAssets() {
	fmt.Printf("Generating assets\n")

	buildFile, err := os.ReadFile(filepath.Join(workspacePath, "configs", "assets.toml"))
	if err != nil {
		panic(err)
	}

	var config struct {
		Images  []GraphicsOutput
		Music   []MusicOutput
		Sounds  []SoundOutput
		Karaoke []KaraokeOutput
	}
	if err := toml.Unmarshal(buildFile, &config); err != nil {
		panic(err)
	}

	assetsPath := filepath.Join(workspacePath, "assets")

	jf := jen.NewFile("assets")

	jf.Comment(warning)
	jf.Line()

	genImagesAssets(jf, config.Images, filepath.Join(assetsPath, "images"))

	// Create sound types
	jf.Type().Id("SoundType").Int()
	jf.Line()
	jf.Const().Id("SoundTypeMp3").Op("SoundType").Op("=").Lit(0)
	jf.Const().Id("SoundTypeWav").Op("SoundType").Op("=").Lit(1)

	genMusicAssets(jf, config.Music, filepath.Join(assetsPath, "music"))
	genSoundAssets(jf, config.Sounds, filepath.Join(assetsPath, "sounds"))
	genKarokeAssets(jf, config.Karaoke, filepath.Join(assetsPath, "karaoke"))

	f, err := os.Create(filepath.Join(workspacePath, "assets", "gen.go"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := jf.Render(f); err != nil {
		panic(err.Error()[:100])
	}
}

func needBuild(strAry []string, value string) bool {
	if len(strAry) == 0 {
		return true
	}

	for _, str := range strAry {
		if strings.EqualFold(str, value) {
			return true
		}
	}

	return false
}

func main() {
	workspacePath = os.Args[1]

	var toBuild []string

	if len(os.Args) > 2 {
		toBuild = os.Args[2:]
	}

	if needBuild(toBuild, "components") {
		genComCode()
	}

	if needBuild(toBuild, "assets") {
		genAssets()
	}
}
