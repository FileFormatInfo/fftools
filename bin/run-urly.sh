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

go build -o ./dist/urly ./cmd/urly

./dist/urly --username=me https://www.example.com/
PASSWORD=thepassword ./dist/urly --password-env=PASSWORD https://me@www.example.com/
