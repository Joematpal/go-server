module github.com/digital-dream-labs/go-server/example/server

go 1.16

require (
	github.com/digital-dream-labs/go-logger v0.0.0-20210716200444-aacb0c191ae0
	github.com/digital-dream-labs/go-server v0.0.0-20210628191752-bf6767d57292
	github.com/urfave/cli/v2 v2.3.0
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/zap v1.19.0 // indirect
	google.golang.org/grpc v1.40.0
)

replace github.com/digital-dream-labs/go-server => ../../
