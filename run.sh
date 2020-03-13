#!/bin/sh

cd "${0%/*}"

buildInfo="`date -u '+%Y-%m-%dT%TZ'`|`git describe --always --long`|`git tag | tail -1`"
go run -ldflags "-X main.buildInfo=${buildInfo} -s -w" ./cmd/sudoku/main.go "$@"
