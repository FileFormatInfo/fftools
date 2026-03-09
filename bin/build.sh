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

CLEAN=1
FILES=("$@")
if [ ${#FILES[@]} -eq 0 ]; then
	echo "INFO: no files specified, building all files in ${REPO_DIR}/cmd"
	CLEAN=0
	FILES=($(ls "${REPO_DIR}/cmd"))
fi
LASTMOD=$(date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT=$(git rev-parse --short HEAD)
if [[ $(git diff --stat) != '' ]]; then
	COMMIT="${COMMIT}-dirty"
fi

for f in "${FILES[@]}"; do
	if [ -f "${DIST_DIR}/${f}" ]; then
		if [ ${CLEAN} -eq 1 ]; then
			echo "INFO: removing existing file ${DIST_DIR}/${f}"
			rm "${DIST_DIR}/${f}"
		else
			echo "WARNING: file ${DIST_DIR}/${f} already exists"
			continue
		fi
	fi
	echo "INFO: compiling ${f} to dist directory"
	go build \
		-ldflags "-s -w -X main.LASTMOD=${LASTMOD} -X main.COMMIT=${COMMIT} -X main.BUILDER=build.sh -X main.VERSION=localdev" \
		-o "${DIST_DIR}/${f}" \
		"${REPO_DIR}/cmd/${f}"
done

echo "INFO: complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
