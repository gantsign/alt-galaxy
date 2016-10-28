#!/bin/bash

set -e

echo
echo '#####################'
echo '# Building binaries #'
echo '#####################'
echo

gox -ldflags="-X github.com/gantsign/alt-galaxy/internal/application.Version=${TRAVIS_TAG:-dev} \
              -X github.com/gantsign/alt-galaxy/internal/application.Revision=${TRAVIS_COMMIT:-unknown} \
              -X github.com/gantsign/alt-galaxy/internal/application.Built=$(date --iso-8601=seconds) \
              -s" \
    -os="${OS_TARGETS:-linux darwin windows}" \
    -arch="${ARCH_TARGETS:-amd64}" \
    -output="dist/{{.OS}}/{{.Arch}}/{{.Dir}}"

# Function to map from exectuable path and archive extension to archive path.
to_archive_path() {
    executable_path="$1"
    archive_extension="$2"

    regex='s:(dist)/([^/]+)/([^/]+)/([^\.]*).*:\1/\4_\2_\3:'

    relative_path="$(echo "${executable_path}" | sed -r -e "${regex}")"
    absolute_path="$(readlink -f "${relative_path}")"
    echo "${absolute_path}${archive_extension}"
}

echo
echo '####################################'
echo '# Packaging Linux redistributables #'
echo '####################################'
echo

find dist/linux -type f | while read executable_path; do

    executable_dir="$(dirname "${executable_path}")"
    executable_name="$(basename "${executable_path}")"
    archive_path="$(to_archive_path "${executable_path}" ".tar.xz")"

    tar --create \
        --xz \
        --file="${archive_path}" \
        --directory="${executable_dir}" \
        --owner=0 \
        --group=0 \
        --mode='u=rwx,go=rx' \
        --verbose \
        "${executable_name}"

done

echo
echo '####################################'
echo '# Packaging macOS redistributables #'
echo '####################################'
echo

find dist/darwin -type f | while read executable_path; do

    executable_dir="$(dirname "${executable_path}")"
    executable_name="$(basename "${executable_path}")"
    archive_path="$(to_archive_path "${executable_path}" ".tar.gz")"

    tar --create \
        --gzip \
        --file="${archive_path}" \
        --directory="${executable_dir}" \
        --owner=0 \
        --group=0 \
        --mode='u=rwx,go=rx' \
        --verbose \
        "${executable_name}"

done

echo
echo '######################################'
echo '# Packaging Windows redistributables #'
echo '######################################'
echo

find dist/windows -type f | while read executable_path; do

    executable_dir="$(dirname "${executable_path}")"
    executable_name="$(basename "${executable_path}")"
    archive_path="$(to_archive_path "${executable_path}" ".7z")"

    (cd "${executable_dir}" && 7z a -bd "${archive_path}" "${executable_name}")

done
