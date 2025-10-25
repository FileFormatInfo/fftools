package main

import (
	"bufio"
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha3"
	_ "crypto/sha512"
	"hash"
	"io"
	"regexp"
	"strings"

	"fmt"
	"os"

	_ "golang.org/x/crypto/blake2b"
	_ "golang.org/x/crypto/blake2s"
	_ "golang.org/x/crypto/md4"
	_ "golang.org/x/crypto/ripemd160"

	"github.com/spf13/pflag"
)

// struct with name and hash
type Hasher struct {
	name string
	hash hash.Hash
}

// map of name to Hasher
var hasherMap = map[string]Hasher{}

func printHashes(fileName string, hashers []Hasher) error {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Read the file in chunks and update the hashes
	buf := make([]byte, 64*1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			for _, h := range hashers {
				h.hash.Write(buf[:n])
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	// Print the results
	for _, h := range hashers {
		if len(hashers) > 1 {
			fmt.Printf("%x: %s %s\n", h.hash.Sum(nil), h.name, fileName)
		} else {
			fmt.Printf("%x: %s\n", h.hash.Sum(nil), fileName)
		}
	}

	return nil
}

var nonAlphaNumeric = regexp.MustCompile(`[^a-zA-Z0-9]`)

func main() {
	var hashNames []string
	var list bool
	pflag.StringArrayVar(&hashNames, "hash", []string{}, "Hash algorithms to use (e.g., sha256, sha512)")
	pflag.BoolVar(&list, "list", false, "List available hash algorithms")
	//LATER: output format

	pflag.Parse()

	if list {
		for i := crypto.MD4; i <= crypto.BLAKE2b_512; i++ {
			hashName := nonAlphaNumeric.ReplaceAllString(strings.ToLower(i.String()), "")
			if i.Available() {
				fmt.Fprintf(os.Stdout, "%s: %s\n", hashName, "true")
			} else {
				fmt.Fprintf(os.Stdout, "%s: %s\n", hashName, "false")
			}
		}
		return
	}

	// build hashers map
	for i := crypto.MD4; i <= crypto.BLAKE2b_512; i++ {
		if i.Available() {
			hashName := nonAlphaNumeric.ReplaceAllString(strings.ToLower(i.String()), "")
			theHash := Hasher{name: hashName, hash: i.New()}
			hasherMap[hashName] = theHash
			hasherMap[strings.ToLower(i.String())] = theHash
		}
	}

	var hashers = []Hasher{}
	if len(hashNames) == 0 {
		hashers = append(hashers, hasherMap["sha256"])
	} else {
		for _, name := range hashNames {
			var h = hasherMap[name]
			if h.hash == nil {
				fmt.Fprintf(os.Stderr, "Unknown or unavailable hash algorithm: %s\n", name)
				os.Exit(1)
			}
			hashers = append(hashers, h)
		}
	}

	args := pflag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}
	for _, arg := range args {
		if arg == "-" {
			arg = "/dev/stdin"
		}

		if err := printHashes(arg, hashers); err != nil {
			panic(err)
		}
	}

}
