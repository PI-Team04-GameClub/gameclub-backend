#!/bin/bash
set -e

go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out
