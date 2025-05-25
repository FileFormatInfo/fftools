#!/usr/bin/env bash
#
# run locally
#

set -o errexit
set -o nounset
set -o pipefail

go build -o ./dist/certinfo ./cmd/certinfo

./dist/certinfo https://www.fileformat.info/
./dist/certinfo tmp/cert.pem
./dist/certinfo tmp/cert.der
