.PHONY: example example_client example_grpc


example_client:
	LOG_ENV=dev \
	LOG_ENCODING=console \
	LOG_LEVEL=fatal \
	GRPC_HOST=0.0.0.0 \
	GPRC_PORT=9000 \
	go run example/client/main.go

example_grpc:
	cd example/grpc && \
	LOG_ENV=dev \
	LOG_ENCODING=console \
	LOG_LEVEL=fatal \
	GRPC_HOST=0.0.0.0 \
	GPRC_PORT=9000 \
	go run main.go

example: example_client example_grpc