package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/dave/jennifer/jen"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancoleman/strcase"

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
	Name        string
	Type        string
	Parts       []string
	ScaleFactor int
	FrameWidth  int
	File        string
}

func (g *GraphicsOutput) genImageAsset(jf *jen.File, img image.Image) {
	var fields []jen.Code
	fields = append(fields, jen.Id("Data").String())
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

	jf.Var().Id(name).Op("=").Struct(fields...).BlockFunc(func(jf *jen.Group) {
		jf.Id("Data").Op(":").Lit(string(imgBytes)).Op(",")

		if g.FrameWidth > 0 {
			jf.Id("FrameWidth").Op(":").Lit(g.FrameWidth).Op(",")
		}

		jf.Id("Compressed").Op(":").Lit(useCompressed).Op(",")
	})

	jf.Line()
}

func (g *GraphicsOutput) genTileSetImageAsset(jf *jen.File, img image.Image) {

	ebitenImg := ebiten.NewImageFromImage(img)
	defer ebitenImg.Dispose()

	for i, part := range g.Parts {
		inG := &GraphicsOutput{
			Name: g.Name + part,
			Type: "single",
		}

		subRect := image.Rectangle{
			Min: image.Point{X: i * g.FrameWidth, Y: 0},
			Max: image.Point{},
		}
		subImg := ebitenImg.SubImage(subRect)

		inG.genImageAsset(jf, subImg)
	}
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

type SoundOutput struct {
	Name string
	File string
}

func (s *SoundOutput) genSouundAssetFromFile(jf *jen.File, path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	name = strcase.ToCamel(name)
	name = "Music" + name

	genByteArray(jf, data, name)
}

func genMusicAssets(jf *jen.File, sounds []SoundOutput, assetsPath string) {
	jf.Const().Id("Mp3SampleRate").Op("=").Lit(48000)

	for _, target := range sounds {
		target.genSouundAssetFromFile(jf, filepath.Join(assetsPath, target.File))
	}
}

func genAssets() {
	fmt.Printf("Generating assets\n")

	buildFile, err := os.ReadFile(filepath.Join(workspacePath, "configs", "assets.toml"))
	if err != nil {
		panic(err)
	}

	var config struct {
		Images []GraphicsOutput
		Music  []SoundOutput
	}
	if err := toml.Unmarshal(buildFile, &config); err != nil {
		panic(err)
	}

	assetsPath := filepath.Join(workspacePath, "assets")

	jf := jen.NewFile("assets")

	jf.Comment(warning)
	jf.Line()

	genImagesAssets(jf, config.Images, filepath.Join(assetsPath, "images"))
	genMusicAssets(jf, config.Music, filepath.Join(assetsPath, "music"))

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
