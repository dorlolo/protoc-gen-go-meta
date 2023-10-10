package example

import (
	"github.com/dorlolo/protoc-gen-go-meta/meta"
	"github.com/golang/protobuf/proto"
	"strconv"
)

func (x TestDataType) Value() string {
	extOpts, err := proto.GetExtension(proto.MessageV1(x.Descriptor().Values().ByNumber(x.Number()).Options()), meta.E_EnumValue)
	if err != nil {
		return strconv.Itoa(int(x))
	}
	enumOptions, ok := extOpts.(*string)
	if ok {
		return *enumOptions
	}
	return strconv.Itoa(int(x))
}
