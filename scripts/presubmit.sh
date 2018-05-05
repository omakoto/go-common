#!/bin/bash

set -e

cd "${0%/*}/.."

NOT_FORMATTED=$(gofmt -s -d $(find . -type f -name '*.go'))

if [[ -n "$NOT_FORMATTED" ]] ; then
    echo $'Wrong format:\n'"$NOT_FORMATTED"
    exit 1
fi

go test -v -race ./...                   # Run all the tests with the race detector enabled

echo "Running extra checks..."
go vet ./...                             # go vet is the official Go static analyzer
megacheck ./...                          # "go vet on steroids" + linter
golint -set_exit_status src/... | grep -v -P '\bexported (type|func|function|method|const|var)\b'
