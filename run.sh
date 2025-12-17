#!/usr/bin/env bash
#
# run tests locally
#

set -o errexit
set -o pipefail
set -o nounset

if [ -d "./dist" ]; then
  echo "INFO: removing old binaries"
  rm -rf ./dist
fi

mkdir ./dist

echo "INFO: building new binaries"
go build -o ./dist/urly cmd/urly/urly.go
if [ ! -f "./dist/urly" ]; then
  echo "ERROR: failed to build urly"
  exit 1
fi


export PATH=$(pwd)/dist:$PATH

echo "INFO: running tests"
go test -timeout 30s -run "^TestUrly$" github.com/FileFormatInfo/fftools/cmd/urly

echo "INFO: complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
