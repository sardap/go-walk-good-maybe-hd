package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"

	_ "github.com/oov/psd"
)

const warning = "AUTO GENERATED CODE DO NOT EDIT REFER TO gen/codegen"

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
			if err != nil || info.IsDir() || info.Name() == "gen.go" {
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

type GraphicsOutput struct {
	Name        string
	ScaleFactor int
	Files       []string
	Dirs        []string
}

func genAssetFile(jf *jen.File, path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	name := strings.TrimSuffix(filepath.Base(f.Name()), filepath.Ext(f.Name()))
	name = strcase.ToCamel(name)
	name = "Image" + name

	buffer := &bytes.Buffer{}
	err = png.Encode(buffer, img)
	if err != nil {
		panic(err)
	}

	compressed := &bytes.Buffer{}
	func() {
		zw := gzip.NewWriter(compressed)
		defer zw.Close()

		_, err = zw.Write(buffer.Bytes())
		if err != nil {
			panic(err)
		}
	}()

	fmt.Printf("Reduced %d by %d bytes\n", buffer.Len(), compressed.Len()-buffer.Len())

	jf.Var().Id(name).Op("=").Index().Byte().ValuesFunc(func(g *jen.Group) {
		for _, b := range compressed.Bytes() {
			g.LitByte(b)
		}
	})
	jf.Line()
}

func genAssets() {
	fmt.Printf("Generating assets\n")

	buildFile, err := os.ReadFile(filepath.Join(workspacePath, "configs", "assets.toml"))
	if err != nil {
		panic(err)
	}

	var config struct {
		Targets []GraphicsOutput
	}
	if err := toml.Unmarshal(buildFile, &config); err != nil {
		panic(err)
	}

	assetsPath := filepath.Join(workspacePath, "assets")

	jf := jen.NewFile("assets")

	jf.Comment(warning)
	jf.Line()

	for _, target := range config.Targets {
		filepath.Walk(assetsPath, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() || strings.Contains(filepath.Base(info.Name()), ".go") {
				return nil
			}

			target.Files = append(target.Files, strings.TrimPrefix(path, assetsPath+string(filepath.Separator)))
			return nil
		})

		for _, fPath := range target.Files {
			genAssetFile(jf, filepath.Join(assetsPath, fPath))
		}
	}

	f, err := os.Create(filepath.Join(workspacePath, "assets", "gen.go"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	jf.Render(f)
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
