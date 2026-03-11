# File Format Tools [<img alt="Logo for fftools" src="docs/favicon.svg" height="96" align="right"/>](https://www.fileformat.info/)

[![build](https://github.com/FileFormatInfo/fftools/actions/workflows/build.yaml/badge.svg)](https://github.com/FileFormatInfo/fftools/actions/workflows/build.yaml)
[![release](https://github.com/FileFormatInfo/fftools/actions/workflows/release.yaml/badge.svg)](https://github.com/FileFormatInfo/fftools/actions/workflows/release.yaml)

## Installation

```
brew install FileFormatInfo/homebrew-tap/fftools
```

Or download from [Releases](https://github.com/FileFormatInfo/fftools/releases)

## Programs

- [asciify](cmd/asciify/README.md): converts to ASCII using [anyascii](https://github.com/anyascii/anyascii)
- [asciitable](cmd/asciitable/README.md): prints out an table of ASCII characters
- [bytecount](cmd/bytecount/README.md): counts the number of occurrences of each byte
- [certinfo](cmd/certinfo/README.md): print info about an x509 (aka SSL/HTTPS) certificate
- [ghash](cmd/ghash/README.md): calculate file hashes
- [hexdumpc](cmd/hexdumpc/README.md): generate canonical hexdump (`hexdump -C`) output in case you don't have  [`hexdump`](https://man7.org/linux/man-pages/man1/hexdump.1.html) installed
- [hosty](cmd/hosty/README.md): manipulate hostnames
- [jsonfmt](cmd/jsonfmt/README.md): format JSON (expanded, canonical, line, fractured)
- [unhexdump](cmd/unhexdump/README.md): convert `hexdump -c` output back into binary
- [unicount](cmd/unicount/README.md): count Unicode codepoints in a file
- [uniwhat](cmd/uniwhat/README.md): print the names of each Unicode character in a file
- [urly](cmd/urly/README.md): manipulate URLs

## Experiments

- `wombat` - tests terminal screen functions

## Credits

[![Git](https://www.vectorlogo.zone/logos/git-scm/git-scm-ar21.svg)](https://git-scm.com/ "Version control")
[![Github](https://www.vectorlogo.zone/logos/github/github-ar21.svg)](https://github.com/ "Code hosting")
[![golang](https://www.vectorlogo.zone/logos/golang/golang-ar21.svg)](https://golang.org/ "Programming language")
[![svgrepo](https://www.vectorlogo.zone/logos/svgrepo/svgrepo-ar21.svg)](https://www.svgrepo.com/svg/276165/gardening-tools-rake "Icon")

* [goreleaser](https://goreleaser.com/)

## To Do

- [ ] `body`: prints specific lines of a file (in between `head` and `tail`)
- [ ] `bom-defuse`: remove byte-order-marks (BOMs) from files
- [ ] `purify`: remove high bytes | non-UTF8 | non-ASCII | etc
- [ ] `trilobyte`: translates bytes according to a map
- [ ] `trune`: translates Unicode codepoints (runes) according to a map
- [ ] `uncolor`: remove color codes (or all terminal escapes) from stdin (see [unansi.c](https://github.com/414owen/unansi/blob/main/unansi.c))
- [ ] `ustrings`: like the standard [`strings`](https://man7.org/linux/man-pages/man1/strings.1.html) utility, but with Unicode support
- [ ] `utf8ify`: convert to UTF-8 (also see [unormalize](https://github.com/eddieantonio/unormalize))
- [ ] `xmlfmt`: pretty print xml
- [ ] `yamlfmt`: pretty print yaml


## General
