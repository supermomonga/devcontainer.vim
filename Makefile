APP_NAME := devcontainer.vim
GOARCH := amd64
WINDOWS_BINARY_NAME := ${APP_NAME}-windows-${GOARCH}.exe
LINUX_BINARY_NAME := ${APP_NAME}-linux-${GOARCH}
DARWIN_BINARY_NAME := ${APP_NAME}-darwin-${GOARCH}

GO_BIN := ${GOPATH}/bin
LD_FLAGS := "-s -w"
VERSION := 0.0.1

DEST := ./build

WATCH_SRC := ./main.go \
						 ./devcontainer/devcontainer.go \
						 ./dockerRun/dockerRun.go \
						 ./tools/tools.go \
						 ./tools/tools_windows.go \
						 ./tools/tools_nowindows.go \
						 ./util/util.go

### 開発関連
# 開発環境の都合で、個別にビルドできるようにしている
# (Linux コンテナ上でコーディングを行い、 Windows 上で実行することがあるため)
all: build
build: build/devcontainer.vim
build/devcontainer.vim: ${WATCH_SRC}
	go build -ldflags=${LD_FLAGS} -trimpath -o ./build/${APP_NAME} ${SRC}

build-all: build-windows build-linux build-darwin

build-windows: build/${WINDOWS_BINARY_NAME}
build/${WINDOWS_BINARY_NAME}: ${WATCH_SRC}
	GOOS=windows GOARCH=${GOARCH} go build -ldflags=${LD_FLAGS} -trimpath -o build/${WINDOWS_BINARY_NAME}

build-linux: build/${LINUX_BINARY_NAME}
build/${LINUX_BINARY_NAME}: ${WATCH_SRC}
	GOOS=linux GOARCH=${GOARCH} go build -ldflags=${LD_FLAGS} -trimpath -o build/${LINUX_BINARY_NAME}

build-darwin: build/${DARWIN_BINARY_NAME}
build/${DARWIN_BINARY_NAME}: ${WATCH_SRC}
	GOOS=darwin GOARCH=${GOARCH} go build -ldflags=${LD_FLAGS} -trimpath -o build/${DARWIN_BINARY_NAME}

.PHONY: clean
clean:
	rm -rf build

