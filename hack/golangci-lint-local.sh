#!/bin/sh

set -e

golangci-lint \
    run \
    --timeout=5m \
    -E bodyclose \
    -E durationcheck \
    -E goimports \
    -E gosec \
    -E misspell \
    -E nolintlint \
    -E nakedret \
    -E unconvert \
    -E unparam \
    -E prealloc \
    -E whitespace \
    -E exportloopref \
    -E revive \
    -E gofmt \
    cmd/... pkg/...
