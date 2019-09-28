PROJECT?=github.com/rumyantseva/paris

RELEASE?=0.0.0
COMMIT := git-$(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

build:
	CGO_ENABLED=0 go build \
		-ldflags "-s -w -X ${PROJECT}/internal/version.Version=${RELEASE} \
		-X ${PROJECT}/internal/version.Commit=${COMMIT} \
		-X ${PROJECT}/internal/version.BuildTime=${BUILD_TIME}" \
		-o bin/paris ${PROJECT}/cmd/paris
