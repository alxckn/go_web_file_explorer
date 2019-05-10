#!/bin/bash

# docker run --rm golang:1.11-alpine go tool dist list | grep arm
# docker run --rm -v "$PWD":/go/src/go_web_file_explorer -w /go/src/go_web_file_explorer -e GOOS=linux -e GOARCH=arm golang:1.11-alpine go build -v
docker run --rm -v "$PWD":/go/src/go_web_file_explorer -w /go/src/go_web_file_explorer golang:1.11 /bin/sh -c "go build && chown 1000:1000 /go/src/go_web_file_explorer/go_web_file_explorer"
