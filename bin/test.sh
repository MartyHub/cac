#!/usr/bin/env sh

set -eux

test='./...'

if [[ ! -z "$@" ]]; then
  test="-run $@"
fi

chmod 600 internal/.config/cac/config.json

go test -coverprofile coverage.out -race -short -timeout 30s ${test}
