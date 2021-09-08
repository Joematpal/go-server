# go-server


## grpc/
----
*Contains GRPC Server code for setting up a http2 Server and code for create a REST Gateway*

Notes:
1. Examples are in the example folder
2. There are already defined flags for the cli implementation (urfave/v2/cli) in `flags/grpc.go`
3. The server code can be used to setup a gateway to an gRPC service that is on a different host and port
4. A gateway server only implementation can not be created as a one to many different gRPC services.
5. TODO: handle multiple different versions of grpc code
6. You can set a version path it is in `grpc/options.go`