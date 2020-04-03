#!/bin/sh

cd "${0%/*}"

go test -run=XXX -bench=. -v $(go list ./generator/...)
