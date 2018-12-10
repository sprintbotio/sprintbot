SHELL= /bin/bash
TAG=dev
OS=linux


.PHONY: run
run:
	source .env && go run cmd/bot/main.go
.PHONY: build
build:
	-rm ./build/sprintbot
	cd cmd/bot && GOOS=${OS} GOARCH=amd64 CGO_ENABLED=0 go build -o ../../build/sprintbot

.PHONY: build_image
build_image: build
	cd build && docker build -t quay.io/sprintbot/sprintbot:${TAG} .

.PHONY: check-gofmt
check-gofmt:
	diff -u <(echo -n) <(gofmt -d `find . -type f -name '*.go' -not -path "./vendor/*"`)

.PHONY: test-unit
test-unit:
	@echo Running tests:
	go test -v -race -cover ./pkg/...

.PHONY: generate
generate:
	@go generate ./...