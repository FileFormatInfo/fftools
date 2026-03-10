# jsonfmt

Format JSON files either expanded, compact or single-line

Default mode is `--fractured` when no mode flag is provided.

## Options

* `--canonical`: the same output as `jq . --sort-keys`
* `--sort-keys`: sort object keys case-insensitively
* `--line`: everything on a single line
* `--trailing-newline`: if a trailing newline should be emitted
* `--eol`: end-of-line character(s) `[ lf | cr' | crlf ]`.  Default is `lf`.
* `--fractured`: compact output with [FracturedJson](https://github.com/j-brooke/FracturedJson) using the [go-fractured-json](https://github.com/FileFormatInfo/go-fractured-json) library

## Fractured Options

`jsonfmt` exposes `go-fractured-json` options as `--fractured-*` flags, including:

* `--fractured-max-total-line-length`
* `--fractured-max-inline-complexity`
* `--fractured-max-compact-array-complexity`
* `--fractured-max-table-row-complexity`
* `--fractured-max-prop-name-padding`
* `--fractured-colon-before-prop-name-padding`
* `--fractured-table-comma-placement` (`before-padding`, `after-padding`, `before-padding-except-numbers`)
* `--fractured-min-compact-array-row-items`
* `--fractured-always-expand-depth`
* `--fractured-nested-bracket-padding`
* `--fractured-simple-bracket-padding`
* `--fractured-colon-padding`
* `--fractured-comma-padding`
* `--fractured-comment-padding`
* `--fractured-number-list-alignment` (`left`, `right`, `decimal`, `normalize`)
* `--fractured-indent-spaces`
* `--fractured-use-tab-to-indent`
* `--fractured-prefix-string`
* `--fractured-comment-policy` (`error`, `remove`, `preserve`)
* `--fractured-preserve-blank-lines`
* `--fractured-allow-trailing-commas`
* `--fractured-json-eol-style` (`default`, `crlf`, `lf`)
