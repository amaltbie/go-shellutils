
USERSTR := $(shell id -u):$(shell id -g)

.docker:
	docker build -t shellutils-build docker/
	touch .docker

.cache:
	mkdir -p .cache

.gopath:
	mkdir -p .gopath

.build:
	mkdir -p .build

.build/ls-darwin: $(SOURCE) .docker | .cache .gopath .build
	docker run \
		-v $(PWD):/build \
		-u $(USERSTR) \
		-v $(PWD)/.cache:/.cache \
		-v $(PWD)/.gopath:/gopath \
		-e GOOS=darwin \
		-e GO111MODULE=on \
		-e GOPATH=/gopath \
		-w /build shellutils-build \
		go build \
		-o .build/ls-darwin \
		cmd/ls.go

.build/ls-linux: $(SOURCE) .docker | .cache .gopath .build
	docker run \
		-v $(PWD):/build \
		-u $(USERSTR) \
		-v $(PWD)/.cache:/.cache \
		-v $(PWD)/.gopath:/gopath \
		-e GOOS=linux \
		-e GO111MODULE=on \
		-e GOPATH=/gopath \
		-w /build shellutils-build \
		go build \
		-o .build/ls-linux \
		cmd/ls.go

build: .build/ls-darwin .build/ls-linux

clean:
	rm -rf .build

clobber: clean
	rm -f .docker
	rm -rf .gopath
	rm -rf .cache

.DEFAULT_GOAL := build
