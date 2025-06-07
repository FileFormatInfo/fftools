#!/usr/bin/env bash
#
# run locally
#

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

cd "${REPO_DIR}"

go build -o ./dist/certinfo ./cmd/certinfo

./dist/certinfo https://www.fileformat.info/
./dist/certinfo tmp/cert.pem
./dist/certinfo tmp/cert.der
