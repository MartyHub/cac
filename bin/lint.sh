#!/usr/bin/env sh

set -eux

go vet
golangci-lint run
