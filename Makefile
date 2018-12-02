SHELL= /bin/bash
TAG=dev


.PHONY: run
run:
	source .env && go run cmd/bot/main.go
.PHONY: build
build:
	cd cmd/bot && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../build/sprintbot

.PHONY: build_image
build_image: build
	cd build && docker build -t quay.io/sprintbot/sprintbot:${TAG} .