#!/usr/bin/env bash
#
# build binaries
#

set -o errexit
set -o pipefail
set -o nounset

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

MANPAGE_DIR="${REPO_DIR}/manpages"

echo "INFO: starting at $(date -u +%Y-%m-%dT%H:%M:%SZ)"

if [ ! -d "${MANPAGE_DIR}" ]; then
	echo "INFO: creating manpage directory ${MANPAGE_DIR}"
	mkdir "${MANPAGE_DIR}"
else
	echo "INFO: using existing manpage directory ${MANPAGE_DIR}"
fi

FILES=$(ls "${REPO_DIR}/cmd")

for f in $FILES; do
	if [ -f "${MANPAGE_DIR}/${f}" ]; then
		echo "WARNING: file ${MANPAGE_DIR}/${f} already exists"
		continue
	fi
	echo "INFO: compiling ${f}"
	#LATER: go build -o "${MANPAGE_DIR}/${f}" "${REPO_DIR}/cmd/${f}"
done

echo "INFO: complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
