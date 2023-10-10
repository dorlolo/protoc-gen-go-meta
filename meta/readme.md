## 插件功能说明
为枚举类型添加Value方法来输出自定义值。比如定义了如下内容
```go
type AuditType int32

const (
    AuditType_pending AuditType = 0
)
```
那么使用 `AuditType_pending.Value()` 可以输出"待审核"。只需要在proto文件中定义好参数，就不需要额外写转换的方法了

## 使用方法
### 1. 准备工作
在项目根目录下，使用以下命令生成protoc-gen-go-meta.exe，然后添加到环境变量
```cmd
go build -o protoc-gen-go-meta.exe main.go
```

复制 [.\meta\metaproto\meta.proto](test%2Ftest.proto) 到你项目的pb目录下

### 2.使用 meta.proto 配置枚举值
导入meta.proto，使用参数`meta.enum_value`添加自定义值（注意两侧要有括号）。
```protobuf
syntax = "proto3";
package test;
option go_package = "./;test";

import "metaproto/meta.proto";//meta.proto的相对位置


enum TestDataType {
  One = 0 [(meta.enum_value) = "一"];
  Two = 1 [(meta.enum_value) = "二"];
}
```

### 3. 运行protoc命令生成go文件
命令如下，也就是多加一个参数`--go-meta_out=` 指定输出路径即可
```cmd
protoc --proto_path=.\meta --go-meta_out=.\meta\test --go_out=.\meta\test .\meta\test\test.proto
```

如果没有将插件添加到环境变量，可参考下面的命令来指定插件路径
```cmd
protoc --plugin=protoc-gen-go-meta=C:\Users\lenovo\go\bin\protoc-gen-go-meta.exe  --proto_path=.\meta --go-meta_out=.\meta\test --go_out=.\meta\test .\meta\test\test.proto
```

3. 调用Vlaues()方法
具体示例请查看[test.pb.enumValue_test.go](test%2Ftest.pb.enumValue_test.go)


## 开发测试
1. 编译到protoc-gen-go-meta目录下
在项目根目录执行此命令
```cmd
go build -o protoc-gen-go-meta.exe ./meta/main.go
```

2. 将./protoc-gen-go-meta/proto/meta.proto 复制到 ./protoc-gen-go-meta/test目录下，与test.proto同级

3. 运行protoc命令，将生成protoc-gen-go-meta/test/test.pb.meta.go文件
```cmd
protoc --plugin=protoc-gen-go-meta=.\protoc-gen-go-meta.exe  --proto_path=.\meta --go_out=.\meta\test --go-meta_out=.\meta\test .\meta\test\test.proto
```