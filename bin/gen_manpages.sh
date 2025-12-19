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

PROGRAMS=$(ls "${REPO_DIR}/cmd")

for PROGRAM in $PROGRAMS; do
	if [ -f "${MANPAGE_DIR}/${PROGRAM}.1" ]; then
		echo "WARNING: file ${MANPAGE_DIR}/${PROGRAM}.1 already exists"
		continue
	fi
	echo "INFO: compiling ${PROGRAM} manpage to ${MANPAGE_DIR}/${PROGRAM}.1"
	if [ ! -f "${REPO_DIR}/cmd/${PROGRAM}/README.md" ]; then
		echo "WARNING: missing README.md for ${PROGRAM}, skipping"
		continue
	fi
	# generate manpage from README.md
	pandoc --standalone --to man "${REPO_DIR}/cmd/${PROGRAM}/README.md" -o "${MANPAGE_DIR}/${PROGRAM}.1"
done

echo "INFO: complete at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
