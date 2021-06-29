.PHONY: example example_client example_grpc


example_client:
	cd example/client && \
	LOG_ENV=dev \
	LOG_ENCODING=console \
	LOG_LEVEL=fatal \
	GRPC_HOST=0.0.0.0 \
	GRPC_PORT=9000 \
	go run main.go

example_grpc:
	cd example/grpc && \
	LOG_ENV=dev \
	LOG_ENCODING=console \
	LOG_LEVEL=fatal \
	GRPC_HOST=0.0.0.0 \
	GRPC_PORT=9000 \
	go run main.go

example: example_client example_grpc