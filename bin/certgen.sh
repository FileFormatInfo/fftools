#!/usr/bin/env bash
#
# generate a self-signed certificate with OpenSSL for testing certinfo
#

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

mkdir -p "${REPO_DIR}/tmp"

openssl req \
	-x509 \
	-newkey rsa:4096 \
	-keyout "${REPO_DIR}/tmp/key.pem" \
	-out "${REPO_DIR}/tmp/cert.pem" \
	-sha256 \
	-days 365 \
	-nodes \
	-subj "/C=US/ST=Pennsylvania/L=Philadelphia/O=ExampleCo/CN=example.com"
