NAME=protoc-gateway-impl

.PHONY: build
build: .compile

.compile:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 \
		go build \
			-mod vendor \
			-ldflags="-w -s" \
			-o bin/${NAME} \
			./cmd/main.go

.PHONY: example
example: .compile .example_gen .example_run

.example_gen:
	@(cd example && buf generate)
.example_run:
	@(cd example && go run cmd/main.go)
