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
	"regexp"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/mjibson/go-dsp/wav"
	"github.com/sardap/walk-good-maybe-hd/common"
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

type KaraokeSession struct {
	Inputs           []KaraokeInput      `json:"inputs"`
	Backgrounds      []KaraokeBackground `json:"backgrounds"`
	SoundFiles       map[string]string   `json:"sounds_files,omitempty"`
	MusicFile        string              `json:"music_file,omitempty"`
	TextImageFiles   map[string]string   `json:"text_image_files"`
	TitleScreenImage string              `json:"title_screen_file"`
}

func encodeJpeg(imagePath string) string {
	f, err := os.Open(imagePath)
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
}

func encodePng(imagePath string) string {
	f, err := os.Open(imagePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	buffer := &bytes.Buffer{}
	if err := png.Encode(buffer, img); err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

func (s *KaraokeOutput) genKaraokeAssetFromFile(jf *jen.File, path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	karaokeIn := KaraokeSession{}
	err = json.Unmarshal(data, &karaokeIn)
	if err != nil {
		panic(err)
	}

	karaokeOut := common.KaraokeSession{
		Sounds:     make(map[string]string),
		TextImages: make(map[string]string),
	}

	karaokePath := strings.TrimSuffix(path, filepath.Base(path))

	karaokeOut.TitleScreenImage = encodeJpeg(filepath.Join(karaokePath, karaokeIn.TitleScreenImage))

	for _, background := range karaokeIn.Backgrounds {
		imageRaw := encodeJpeg(filepath.Join(karaokePath, background.Image))
		karaokeOut.Backgrounds = append(karaokeOut.Backgrounds,
			&common.KaraokeBackground{
				Duration: time.Duration(background.Duration) * time.Millisecond,
				FadeIn:   time.Duration(background.FadeIn) * time.Millisecond,
				Image:    imageRaw,
			},
		)
	}

	for _, input := range karaokeIn.Inputs {
		karaokeOut.Inputs = append(karaokeOut.Inputs, &common.KaraokeInput{
			StartTime: time.Duration(input.StartTime) * time.Millisecond,
			Duration:  time.Duration(input.Duration) * time.Millisecond,
			Sound:     input.Sound,
		})
	}

	for key, textImageFile := range karaokeIn.TextImageFiles {
		karaokeOut.TextImages[key] = encodePng(filepath.Join(karaokePath, textImageFile))
	}

	// timeElapsed := 0
	// for i := 1; i < len(karaokeIn.Inputs); i++ {
	// 	lastInput := &karaokeIn.Inputs[i-1]
	// 	input := &karaokeIn.Inputs[i]

	// 	if input.StartTime > 0 {
	// 		continue
	// 	}

	// 	input.StartTime = timeElapsed + lastInput.Duration + input.StartOffset
	// 	input.StartOffset = 0
	// 	timeElapsed += input.Duration
	// }

	for key, sound := range karaokeIn.SoundFiles {
		soundRaw, err := ioutil.ReadFile(filepath.Join(karaokePath, sound))
		if err != nil {
			panic(err)
		}
		karaokeOut.Sounds[key] = base64.RawStdEncoding.EncodeToString(soundRaw)
	}

	karaokeOut.SampleRate = getMp3SampleRate(filepath.Join(karaokePath, karaokeIn.MusicFile))
	rawMusic, _ := ioutil.ReadFile(filepath.Join(karaokePath, karaokeIn.MusicFile))
	karaokeOut.Music = base64.RawStdEncoding.EncodeToString(rawMusic)

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(karaokeOut); err != nil {
		panic(err)
	}

	compressed, _ := compressAsset(buf.Bytes())

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	name = strcase.ToCamel(name)
	name = "Karaoke" + name

	fields := []jen.Code{
		jen.Id("JsonStr").String(),
	}

	jf.Var().Id(name).Op("=").Struct(fields...).BlockFunc(func(jf *jen.Group) {
		jf.Id("JsonStr").Op(":").Lit(string(compressed)).Op(",")
	})
}

func genKarokeAssets(jf *jen.File, karaoke []KaraokeOutput, assetsPath string) {
	fmt.Printf("\nGenerating Karaoke\n")

	for _, target := range karaoke {
		fmt.Printf("...%s\n", target.Name)
		target.genKaraokeAssetFromFile(jf, filepath.Join(assetsPath, target.File))
	}
}

type IconOutput struct {
	Name       string
	Path       string
	Expression string
	Type       string
	Buttons    GamepadButtonMap
}

type GamepadButtonMap struct {
	A          string `toml:"a"`
	B          string `toml:"b"`
	X          string `toml:"x"`
	Y          string `toml:"y"`
	LB         string `toml:"lb"`
	RB         string `toml:"rb"`
	Select     string `toml:"select"`
	Start      string `toml:"start"`
	LeftStick  string `toml:"left_stick"`
	RightStick string `toml:"right_stick"`
}

func findGroupIdx(key string, keys []string) int {
	result := -1
	for i := 1; i < len(keys); i++ {
		if keys[i] == key {
			result = i * 2
			break
		}
	}

	return result
}

func (i *IconOutput) genKeyboardIconAsset(jf *jen.File, path string) {
	re := regexp.MustCompile(i.Expression)
	keys := re.SubexpNames()
	nameGroup := findGroupIdx("name", keys)

	inputMap := make(map[string][]byte)

	items, _ := ioutil.ReadDir(path)
	for _, item := range items {
		if item.IsDir() {
			continue
		}

		matches := re.FindAllStringSubmatchIndex(item.Name(), -1)
		if len(matches) == 0 {
			continue
		}
		rawName := item.Name()[matches[0][nameGroup]:matches[0][nameGroup+1]]

		data, err := ioutil.ReadFile(filepath.Join(path, item.Name()))
		if err != nil {
			panic(err)
		}

		cleanedName := "Key" + strcase.ToCamel(strings.ReplaceAll(rawName, "_", " "))
		inputMap[cleanedName] = data
	}

	name := "Icon" + strcase.ToCamel(i.Name)

	jf.Var().Id(name).Op("=").Map(jen.Op("ebiten.Key")).String().Values(jen.DictFunc(func(d jen.Dict) {
		for key, image := range inputMap {
			d[jen.Qual("github.com/hajimehoshi/ebiten/v2", key)] = jen.Lit(string(image))
		}
	}))
	jf.Line()
}

func (i *IconOutput) genGamepadIconAsset(jf *jen.File, path string) {
	inputMap := make(jen.Dict)

	loadData := func(fileName string) string {
		data, err := ioutil.ReadFile(filepath.Join(path, fileName))
		if err != nil {
			panic(err)
		}
		return string(data)
	}

	inputMap[jen.Lit(0)] = jen.Lit(loadData(i.Buttons.A))
	inputMap[jen.Lit(1)] = jen.Lit(loadData(i.Buttons.B))
	inputMap[jen.Lit(2)] = jen.Lit(loadData(i.Buttons.X))
	inputMap[jen.Lit(3)] = jen.Lit(loadData(i.Buttons.Y))
	inputMap[jen.Lit(4)] = jen.Lit(loadData(i.Buttons.LB))
	inputMap[jen.Lit(5)] = jen.Lit(loadData(i.Buttons.RB))
	inputMap[jen.Lit(6)] = jen.Lit(loadData(i.Buttons.Select))
	inputMap[jen.Lit(7)] = jen.Lit(loadData(i.Buttons.Start))
	inputMap[jen.Lit(8)] = jen.Lit(loadData(i.Buttons.LeftStick))
	inputMap[jen.Lit(9)] = jen.Lit(loadData(i.Buttons.RightStick))

	name := "Icon" + strcase.ToCamel(i.Name)

	jf.Var().Id(name).Op("=").Map(jen.Int()).String().Values(inputMap)
	jf.Line()
}

func genIconAssets(jf *jen.File, icons []IconOutput, assetsPath string) {
	fmt.Printf("\nGenerating Icons\n")

	for _, target := range icons {
		fmt.Printf("...%s\n", target.Name)
		switch target.Type {
		case "keyboard":
			target.genKeyboardIconAsset(jf, filepath.Join(assetsPath, target.Path))
		case "gamepad":
			target.genGamepadIconAsset(jf, filepath.Join(assetsPath, target.Path))
		}
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
		Icons   []IconOutput
	}
	if err := toml.Unmarshal(buildFile, &config); err != nil {
		panic(err)
	}

	assetsPath := filepath.Join(workspacePath, "assets")

	jf := jen.NewFile("assets")

	jf.ImportName("github.com/hajimehoshi/ebiten/v2", "ebiten")

	jf.Comment(warning)
	jf.Line()

	genImagesAssets(jf, config.Images, filepath.Join(assetsPath, "images"))

	genMusicAssets(jf, config.Music, filepath.Join(assetsPath, "music"))
	genSoundAssets(jf, config.Sounds, filepath.Join(assetsPath, "sounds"))
	genKarokeAssets(jf, config.Karaoke, filepath.Join(assetsPath, "karaoke"))

	genIconAssets(jf, config.Icons, filepath.Join(assetsPath, "inputIcons"))

	f, err := os.Create(filepath.Join(workspacePath, "assets", "gen.go"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := jf.Render(f); err != nil {
		f.Close()
		os.Remove(filepath.Join(workspacePath, "assets", "gen.go"))
		panic(err.Error()[:5000])
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
