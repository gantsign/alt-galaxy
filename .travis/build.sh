#!/bin/bash

set -e

echo
echo '#####################'
echo '# Building binaries #'
echo '#####################'
echo

gox -ldflags="-X main.version=${TRAVIS_TAG:-unknown} \
              -X main.revision=${TRAVIS_COMMIT:-unknown} \
              -X main.built=$(date --iso-8601=seconds) \
              -s" \
    -os="${OS_TARGETS:-linux darwin windows}" \
    -arch="${ARCH_TARGETS:-amd64}" \
    -output="dist/{{.OS}}/{{.Arch}}/{{.Dir}}"
