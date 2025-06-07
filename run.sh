#!/usr/bin/env bash
#
# run locally
#

set -o errexit
set -o nounset
set -o pipefail

go build -o ./dist/asciitable ./cmd/asciitable

./dist/asciitable
