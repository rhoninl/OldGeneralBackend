PROJECT_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

build-image:
	docker buildx build --platform=linux/$(shell go env GOARCH) -f ${PROJECT_ROOT}/dockerfiles/dockerfile \
	--build-arg PROJECT_ROOT="${PROJECT_ROOT}" ${PROJECT_ROOT} \
	-t imagenage:version --load