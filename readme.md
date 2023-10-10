## Description
Adds the Values() method to the enum class so that it can output custom string value.

## How to use ?
1. make sure your/go/bin/path is in your environment variables.

2. install the plugin.
```cmd
go install github.com/dorlolo/protoc-gen-go-meta
```

3. copy [./mea.proto](./mea.proto) into your project.

4. then, You can refer to [./example/test.proto](./example/test.proto) for development 

5. finally, using protoc to produce the file, just add a '--go-meta_out' parameter to the original command. 

for example:
    ```cmd
    protoc --proto_path=. --go-meta_out=.\example --go_out=.\example .\example\test.proto
    ```