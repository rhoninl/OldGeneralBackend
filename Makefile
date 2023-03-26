PROJECT_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

build-image: build-image-api build-image-iam build-image-user build-image-flags build-image-wallet

build-image-api:
	docker buildx build --platform=linux/amd64 \
	-f ${PROJECT_ROOT}/dockerfiles/dockerfile.api . \
	-t serviceapi:nightly --load

build-image-iam:
	docker buildx build --platform=linux/amd64 \
	-f ${PROJECT_ROOT}/dockerfiles/dockerfile.iam . \
	-t serviceiam:nightly --load

build-image-user:
	docker buildx build --platform=linux/amd64 \
	-f ${PROJECT_ROOT}/dockerfiles/dockerfile.user . \
	-t serviceuser:nightly --load

build-image-flags:
	docker buildx build --platform=linux/amd64 \
	-f ${PROJECT_ROOT}/dockerfiles/dockerfile.flags . \
	-t serviceflags:nightly --load

build-image-wallet:
	docker buildx build --platform=linux/amd64 \
	-f ${PROJECT_ROOT}/dockerfiles/dockerfile.wallet . \
	-t servicewallet:nightly --load

test: fmt
	go test -v -race -coverprofile=coverage.out -covermode=atomic $(shell go list ./...)

build-protos: build-protos-go build-protos-swift

build-protos-go:
	protoc --go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	--experimental_allow_proto3_optional \
	Proto/*/*.proto

build-protos-swift:
	protoc --swift_out=. --grpc-swift_out=Client=true,Server=false:. Proto/*/*.proto

install-dependencies:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/golang/mock/mockgen@v1.6.0

fmt:
	go fmt ./...
