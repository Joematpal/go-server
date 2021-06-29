module github.com/digital-dream-labs/go-server/example/server

go 1.16

require (
	github.com/digital-dream-labs/go-logger v0.0.0-20210629173042-a947afccbe42
	github.com/digital-dream-labs/go-server v0.0.0-20210628191752-bf6767d57292
	github.com/urfave/cli/v2 v2.3.0
	go.uber.org/zap v1.18.1 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20210629170331-7dc0b73dc9fb // indirect
	google.golang.org/genproto v0.0.0-20210629135825-364e77e5a69d // indirect
	google.golang.org/grpc v1.38.0
	google.golang.org/grpc/examples v0.0.0-20210628165121-83f9def5feb3
)

replace github.com/digital-dream-labs/go-server => ../../
