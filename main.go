// Package protoc_gen_go_enumvalue -----------------------------
// @file      : main.go
// @author    : JJXu
// @contact   : wavingbear@163.com
// @time      : 2023/10/7 22:17
// -------------------------------------------
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/dorlolo/protoc-gen-go-meta/meta"
	gengo "google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
)

var version = "dev"

func main() {
	versionFlag := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			GenerateFile(gen, f)
		}
		gen.SupportedFeatures = gengo.SupportedFeatures
		return nil
	})
}

func hasEnumValueOption(file *protogen.File) bool {
	for _, e := range file.Enums {
		for _, v := range e.Values {
			opts := v.Desc.Options()
			if opts != nil && proto.HasExtension(proto.MessageV1(opts), meta.E_EnumValue) {
				return true
			}
		}
	}
	return false
}

func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if !hasEnumValueOption(file) {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + ".pb.enumValue.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("import (")
	g.P("\"github.com/golang/protobuf/proto\"")
	g.P("\"github.com/dorlolo/protoc-gen-go-meta/meta\"")
	g.P("\"strconv\"")
	g.P(")")

	for _, e := range file.Enums {
		g.P("func (x ", e.GoIdent, ") Value() string {")
		g.P("extOpts, err := proto.GetExtension(proto.MessageV1(x.Descriptor().Values().ByNumber(x.Number()).Options()), meta.E_EnumValue)")
		g.P("if err != nil {")
		g.P("return strconv.Itoa(int(x))")
		g.P("}")
		g.P("enumOptions, ok := extOpts.(*string)")
		g.P("if ok {")
		g.P("return *enumOptions")
		g.P("}")
		g.P("return strconv.Itoa(int(x))")
		g.P("}")
		g.P()
	}

	return g
}
