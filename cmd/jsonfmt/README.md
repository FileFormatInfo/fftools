# jsonfmt

Format JSON files either expanded, compact or single-line

## Options

* `--canonical`: the same output as `jq . --sort-keys`
* `--line`: everything on a single line
* `--trailing-newline`: if a trailing newline should be emitted
* `--eol`: end-of-line character(s) `[ lf | cr' | crlf ]`.  Default is `lf`.
* `--fractured`: compact output with [FracturedJson](https://github.com/j-brooke/FracturedJson) using the [go-fractured-json](https://github.com/FileFormatInfo/go-fractured-json) library
