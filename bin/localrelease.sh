#!/usr/bin/env bash
#
# build binaries with GoReleaser
#

set -o errexit
set -o pipefail
set -o nounset

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

echo "INFO: starting at $(date -u +%Y-%m-%dT%H:%M:%SZ)"

cd "${REPO_DIR}"
goreleaser --snapshot --clean release

echo "INFO: complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
