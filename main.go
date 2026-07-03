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
	"runtime/debug"

	gengo "google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var version = "dev"

func init() {
	// 读取 Go 构建信息（go install 自动注入）
	if info, ok := debug.ReadBuildInfo(); ok {
		// 优先使用 go install 带来的版本号（v1.0.0）
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
	}
}

func main() {
	versionFlag := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		ext := findEnumValueExtension(gen)
		if ext == nil {
			return nil
		}
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			GenerateFile(gen, f, ext)
		}
		gen.SupportedFeatures = gengo.SupportedFeatures
		return nil
	})
}

func findEnumValueExtension(gen *protogen.Plugin) *protogen.Extension {
	for _, f := range gen.Files {
		for _, ext := range f.Extensions {
			if ext.Desc.Name() == "enum_value" && ext.Desc.ContainingMessage().FullName() == "google.protobuf.EnumValueOptions" {
				return ext
			}
		}
	}
	return nil
}

func hasExtensionField(b []byte, number protoreflect.FieldNumber) bool {
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		if n < 0 {
			return false
		}
		b = b[n:]
		if num == protowire.Number(number) {
			return true
		}
		n = protowire.ConsumeFieldValue(num, typ, b)
		if n < 0 {
			return false
		}
		b = b[n:]
	}
	return false
}

func enumHasEnumValueOption(enum *protogen.Enum, number protoreflect.FieldNumber) bool {
	for _, v := range enum.Values {
		opts := v.Desc.Options()
		if opts != nil && hasExtensionField(opts.ProtoReflect().GetUnknown(), number) {
			return true
		}
	}
	return false
}

func hasEnumValueOption(file *protogen.File, number protoreflect.FieldNumber) bool {
	for _, e := range file.Enums {
		if enumHasEnumValueOption(e, number) {
			return true
		}
	}
	return false
}

func GenerateFile(gen *protogen.Plugin, file *protogen.File, ext *protogen.Extension) *protogen.GeneratedFile {
	enumValueNumber := ext.Desc.Number()
	if !hasEnumValueOption(file, enumValueNumber) {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + ".pb.enumValue.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("package ", file.GoPackageName)
	g.P()

	// Build the correct extension variable identifier with E_ prefix
	extIdent := protogen.GoIdent{
		GoImportPath: ext.GoIdent.GoImportPath,
		GoName:       "E_" + ext.GoIdent.GoName,
	}

	// Pre-declare that we need these imports (protogen will auto-add them when we use GoIdent)
	// But we also need strconv and proto, which we'll reference by string
	strconvIdent := protogen.GoIdent{GoImportPath: "strconv", GoName: "Itoa"}
	protoIdent := protogen.GoIdent{GoImportPath: "google.golang.org/protobuf/proto", GoName: "GetExtension"}

	for _, e := range file.Enums {
		if !enumHasEnumValueOption(e, enumValueNumber) {
			continue
		}
		g.P("func (x ", e.GoIdent, ") Value() string {")
		g.P("opts := x.Descriptor().Values().ByNumber(x.Number()).Options()")
		g.P("if val := ", protoIdent, "(opts, ", extIdent, "); val != nil {")
		g.P("if str, ok := val.(string); ok {")
		g.P("return str")
		g.P("}")
		g.P("}")
		g.P("return ", strconvIdent, "(int(x))")
		g.P("}")
		g.P()
	}

	return g
}
