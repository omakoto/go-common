#!/bin/bash

set -e

cd "${0%/*}/.."

gofmt -s -d $(find . -type f -name '*.go') |& perl -pe 'END{exit($. > 0 ? 1 : 0)}'

go test -v -race ./...                   # Run all the tests with the race detector enabled

echo "Running extra checks..."
go vet ./...                             # go vet is the official Go static analyzer
# staticcheck ./...
golint $(go list ./...) |& grep -v 'exported .* should \(have\|be\)' | perl -pe 'END{exit($. > 0 ? 1 : 0)}'
