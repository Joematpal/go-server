
run:
	@GRPC_PORT=9090 \
	GRPC_TLS=true \
	GRPC_PUB_CERT=./localhost.pem \
	GRPC_PRIV_CERT=./localhost-key.pem \
	SWAGGER_FILE=../../swagger/streamer.swagger.json \
	LOG_LEVEL=debug \
	go run main.go