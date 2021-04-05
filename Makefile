include version.mk

GO_CMD=CGO_ENABLED=0 GO111MODULE=on go
GO_PLATFORM_FLAGS=GOOS=linux GOARCH=amd64
GOPATH?=$(HOME)/go
COMMIT?=$(shell git rev-parse HEAD)

APP?=tic-tak-toe
REGISTRY_PATH?=belazar13/tic-tak-toe

LDFLAGS:=-s -w -X 'github.com/smallhive/tic-tak-toe/app.Name=${APP}' \
		 -X 'github.com/smallhive/tic-tak-toe/app.Commit=${COMMIT}' \
		 -X 'github.com/smallhive/tic-tak-toe/app.Version=${RELEASE}'

build: clean
	$(GO_PLATFORM_FLAGS) $(GO_CMD) build -ldflags "$(LDFLAGS)" \
	-o ./build/bin/${APP} ./cmd/${APP}

clean:
	@rm -f ./build/bin/${APP}

docker-build:
	docker build -t ${REGISTRY_PATH}:${RELEASE} -t ${REGISTRY_PATH}:latest -f ./.cloud/build/Dockerfile .

docker-push:
	docker push ${REGISTRY_PATH}:${RELEASE}
	docker push ${REGISTRY_PATH}:latest
