#!/usr/bin/env bash
#
# run locally
#

set -o errexit
set -o nounset
set -o pipefail

CMD="${1:-bytecount}"

go build -o "./dist/${CMD}" "./cmd/${CMD}"

cat README.md | ./dist/${CMD}
