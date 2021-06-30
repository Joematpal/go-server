.PHONY: example example_client example_server


example_client:
	cd example/client && \
	LOG_ENV=dev \
	LOG_ENCODING=console \
	LOG_LEVEL=fatal \
	GRPC_HOST=0.0.0.0 \
	GRPC_PORT=9000 \
	go run main.go

example_server:
	cd example/server && \
	LOG_ENV=dev \
	LOG_ENCODING=console \
	LOG_LEVEL=fatal \
	GRPC_HOST=0.0.0.0 \
	GRPC_PORT=9000 \
	go run main.go

example_server_gateway:
	cd example/server && \
	LOG_ENV=dev \
	LOG_ENCODING=console \
	LOG_LEVEL=fatal \
	GRPC_HOST=localhost \
	GRPC_PORT=9000 \
	GRPC_TLS=true \
	GRPC_PUB_CERT=../../certs/localhost.pem \
	GRPC_PRIV_CERT=../../certs/localhost-key.pem \
	go run main.go


.PHONY: proto
proto: go.mod bin/protoc bin/protoc-gen-go bin/protobuf
	@./scripts/gen-proto-stubs

bin/go.mod:
	@echo "// Hey, go mod, keep out!" > bin/go.mod

go.mod:
	@go mod init

.PHONY: verify-proto
verify-proto: proto
	@./scripts/git-diff

bin/protoc: scripts/get-protoc
	@./scripts/get-protoc bin/protoc

bin/protobuf: bin/go.mod scripts/get-protoc-extras
	@./scripts/get-protoc-extras bin/protobuf

bin/protoc-gen-go:
	# @go get -u github.com/golang/protobuf/protoc-gen-go