#!/usr/bin/env bash
#
# build binaries
#

set -o errexit
set -o pipefail
set -o nounset

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

echo "INFO: testing starting at $(date -u +%Y-%m-%dT%H:%M:%SZ)"

FILES=("$@")
if [ ${#FILES[@]} -eq 0 ]; then
	echo "INFO: no files specified, building all files in ${REPO_DIR}/cmd"
	FILES=($(ls "${REPO_DIR}/cmd"))
fi

for f in "${FILES[@]}"; do
	if [ ! -d "${REPO_DIR}/cmd/${f}" ]; then
		echo "ERROR: ${f} not found"
		exit 1
	fi
	echo "INFO: testing ${f}"
	go test "${REPO_DIR}/cmd/${f}"
done

echo "INFO: testing complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
