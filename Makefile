# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=echoserver
BINARY_VER=latest
BINARY_WINDOWS=$(BINARY_NAME)_windows.exe
BINARY_LINUX=$(BINARY_NAME)_linux

IMAGE_REPO=registry.cn-hangzhou.aliyuncs.com/gorux
IMAGE_NAME=$(BINARY_NAME):$(BINARY_VER)

CONTAINER_NAME=$(BINARY_NAME)_inst

build:
	$(GOBUILD) -o bin/$(BINARY_NAME) .

clean:
	$(GOCLEAN)

clean-bin:
	rm -f bin/*.*

build4win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -a -o bin/$(BINARY_WINDOWS) .

build4linux:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -a -o bin/$(BINARY_LINUX) .

docker-build:
	docker build . -t $(IMAGE_NAME)

docker-run:
	docker run -d --name $(CONTAINER_NAME) -p 8181:8181 $(IMAGE_NAME)

docker-stop:
	docker stop $(CONTAINER_NAME) && docker rm -v $(CONTAINER_NAME)

docker-push:
	docker tag  $(IMAGE_NAME) $(IMAGE_REPO)/$(IMAGE_NAME)
	docker push $(IMAGE_REPO)/$(IMAGE_NAME)
