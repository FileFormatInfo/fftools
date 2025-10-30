#!/usr/bin/env bash
#
# build binaries
#

set -o errexit
set -o pipefail
set -o nounset

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

DIST_DIR="${REPO_DIR}/dist"

echo "INFO: starting at $(date -u +%Y-%m-%dT%H:%M:%SZ)"

if [ ! -d "${DIST_DIR}" ]; then
	echo "INFO: creating dist directory ${DIST_DIR}"
	mkdir "${REPO_DIR}/dist"
fi

FILES=$(ls "${REPO_DIR}/cmd")
LASTMOD=$(date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT=$(git rev-parse --short HEAD)
if [[ $(git diff --stat) != '' ]]; then
	COMMIT="${COMMIT}-dirty"
fi

for f in $FILES; do
	if [ -f "${DIST_DIR}/${f}" ]; then
		echo "WARNING: file ${DIST_DIR}/${f} already exists"
		continue
	fi
	echo "INFO: compiling ${f} to dist directory"
	go build \
		-ldflags "-s -w -X github.com/FileFormatInfo/fftools/internal.LASTMOD=${LASTMOD} -X github.com/FileFormatInfo/fftools/internal.COMMIT=${COMMIT} -X github.com/FileFormatInfo/fftools/internal.BUILDER=build.sh" \
		-o "${DIST_DIR}/${f}" \
		"${REPO_DIR}/cmd/${f}"
done

echo "INFO: complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
