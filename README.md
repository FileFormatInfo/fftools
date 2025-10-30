# File Format Tools [<img alt="Logo for fftools" src="docs/favicon.svg" height="96" align="right"/>](https://www.fileformat.info/)

[![build](https://github.com/FileFormatInfo/fftools/actions/workflows/build.yaml/badge.svg)](https://github.com/FileFormatInfo/fftools/actions/workflows/build.yaml)


## Programs

- `asciify`: converts to ASCII using [anyascii](https://github.com/anyascii/anyascii)
- `asciitable`: prints out an table of ASCII characters
- `bytecount`: counts the number of occurrences of each byte
- `certinfo`: print info about an x509 (aka SSL/HTTPS) certificate
- `hexdumpc`: generate canonical hexdump (`hexdump -C`) output in case you don't have [`hexdump`](https://man7.org/linux/man-pages/man1/hexdump.1.html) installed
- `unicount`: count Unicode codepoints in a file
- `uniwhat`: print the names of each Unicode character in a file

## Experiments

- `wombat` - tests terminal screen functions
- `spinner` - terminal spinner for long-running tasks

## Credits

[![Git](https://www.vectorlogo.zone/logos/git-scm/git-scm-ar21.svg)](https://git-scm.com/ "Version control")
[![Github](https://www.vectorlogo.zone/logos/github/github-ar21.svg)](https://github.com/ "Code hosting")
[![golang](https://www.vectorlogo.zone/logos/golang/golang-ar21.svg)](https://golang.org/ "Programming language")
[![svgrepo](https://www.vectorlogo.zone/logos/svgrepo/svgrepo-ar21.svg)](https://www.svgrepo.com/svg/276165/gardening-tools-rake "Icon")

## Scripts To Do

- [ ] `ghash`: calculate various [hashes available in the Go standard library](https://pkg.go.dev/crypto#Hash)
- [ ] `body`: prints specific lines of a file (in between `head` and `tail`)
- [ ] `trilobyte`: translates bytes according to a map
- [ ] `ustrings`: like the standard [`strings`](https://man7.org/linux/man-pages/man1/strings.1.html) utility, but with Unicode support
- [ ] `utf8ify`: convert to UTF-8
- [ ] `unhexdump`: convert the (edited) output of [`hexdump -C`](https://man7.org/linux/man-pages/man1/hexdump.1.html) back to binary


## General
