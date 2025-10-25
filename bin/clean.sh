#!/usr/bin/env bash
#
# clean up build artifacts
#

set -o errexit
set -o pipefail
set -o nounset

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

DIST_DIR="${REPO_DIR}/dist"

echo "INFO: starting at $(date -u +%Y-%m-%dT%H:%M:%SZ)"

if [ -d "${DIST_DIR}" ]; then
	echo "INFO: removing dist directory ${DIST_DIR}"
	rm -rf "${DIST_DIR}"
else
	echo "INFO: dist directory ${DIST_DIR} does not exist, nothing to clean"
fi

echo "INFO: complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
