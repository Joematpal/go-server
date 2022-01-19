module github.com/digital-dream-labs/go-server/example/server

go 1.16

require (
	github.com/digital-dream-labs/go-logger v0.0.0-20210716200444-aacb0c191ae0
	github.com/digital-dream-labs/go-server v0.0.0-20211101175551-2d3980582ea6
	github.com/urfave/cli/v2 v2.3.0
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/zap v1.20.0 // indirect
	golang.org/x/net v0.0.0-20220114011407-0dd24b26b47d // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	google.golang.org/genproto v0.0.0-20220114231437-d2e6a121cae0 // indirect
	google.golang.org/grpc v1.43.0
)

replace github.com/digital-dream-labs/go-server => ../../
