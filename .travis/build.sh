#!/bin/bash

set -e

echo
echo '#####################'
echo '# Building binaries #'
echo '#####################'
echo

gox -os="${OS_TARGETS:-linux darwin windows}" \
    -arch="${ARCH_TARGETS:-amd64}" \
    -output="dist/{{.OS}}/{{.Arch}}/{{.Dir}}"
