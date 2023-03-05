PROJECT_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

build-image:
	docker buildx build --platform=linux/$(shell go env GOARCH) -f ${PROJECT_ROOT}/dockerfiles/dockerfile \
	--build-arg PROJECT_ROOT="${PROJECT_ROOT}" ${PROJECT_ROOT} \
	-t imagenage:version --load

build-image-api:
	docker buildx build --platform=linux/amd64 -f ${PROJECT_ROOT}/dockerfiles/dockerfile.api \
	--build-arg PROJECT_ROOT="${PROJECT_ROOT}" ${PROJECT_ROOT} \
	-t serviceapi:nightly --load

test: fmt
	go test -v -race -coverprofile=coverage.out -covermode=atomic $(shell go list ./...)

install-dependencies:
	make -f Proto/Makefile install-go-dependencies

build-protos:
	make -f Proto/Makefile build-protos

fmt:
	go fmt ./...
