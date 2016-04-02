#!/bin/sh

cd $(dirname "$0")
mkdir -p tmp
go test -coverprofile tmp/cover.out
go tool cover -html=tmp/cover.out -o tmp/cover.html
