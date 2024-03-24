BIN := "./bin/final-project"
DOCKER_IMG="final-project:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build-final-project:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd

run: build-final-project
	$(BIN) -config ./config.json

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

test:
	go test -v -count=1 -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/opt v1.55.2

lint: install-lint-deps
	/opt/golangci-lint run ./...

.PHONY: build-final-project run build-img run-img test lint
