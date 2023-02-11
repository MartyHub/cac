#!/usr/bin/env sh

set -eux

test='./...'

if [[ ! -z "$@" ]]; then
  test="-run $@"
fi

go test -coverprofile coverage.out -short -timeout 30s ${test}
