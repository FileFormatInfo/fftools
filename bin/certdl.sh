#!/usr/bin/env bash
#
# download a certificate chain from a URL
#

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
mkdir -p "${REPO_DIR}/tmp"
URL="${1:-www.fileformat.info}"
CERT_FILE="${REPO_DIR}/tmp/cert.pem"

echo "Downloading certificate chain from ${URL} to ${CERT_FILE}"

openssl s_client -showcerts -connect "${URL}:443" </dev/null 2>/dev/null \
	| awk '/-----BEGIN CERTIFICATE-----/,/-----END CERTIFICATE-----/ { print }' > "${CERT_FILE}"
