package example

import (
	strconv "strconv"

	meta "github.com/dorlolo/protoc-gen-go-meta/meta"
	proto "google.golang.org/protobuf/proto"
)

func (x TestDataType) Value() string {
	opts := x.Descriptor().Values().ByNumber(x.Number()).Options()
	if val := proto.GetExtension(opts, meta.E_EnumValue); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return strconv.Itoa(int(x))
}
