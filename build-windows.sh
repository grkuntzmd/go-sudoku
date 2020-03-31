#!/bin/sh

cd "${0%/*}"
# INLINE="-a -gcflags -m"

buildInfo="`date -u '+%Y-%m-%dT%TZ'`|`git describe --always --long`|`git tag | tail -1`"
GOOS=windows GOARCH=amd64 go build ${INLINE} -ldflags "-X main.buildInfo=${buildInfo} -s -w" ./cmd/... #  ./generator/...
