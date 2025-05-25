package main

import (
	"bufio"
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha3"
	_ "crypto/sha512"

	"fmt"
	"os"

	_ "golang.org/x/crypto/blake2b"
	_ "golang.org/x/crypto/blake2s"
	_ "golang.org/x/crypto/md4"
	_ "golang.org/x/crypto/ripemd160"

	"github.com/spf13/pflag"
)

func printHashes(fileName string, hashes []string) error {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	return nil
}

func main() {
	var hashes []string
	var list bool
	pflag.StringArrayVar(&hashes, "hash", []string{}, "Hash algorithms to use (e.g., sha256, sha512)")
	pflag.BoolVar(&list, "list", false, "List available hash algorithms")
	//LATER: output format

	pflag.Parse()

	if list {
		for i := crypto.MD4; i <= crypto.BLAKE2b_512; i++ {
			if i.Available() {
				fmt.Fprintf(os.Stdout, "%s: %s\n", i, "true")
			} else {
				fmt.Fprintf(os.Stdout, "%s: %s\n", i, "false")
			}
		}
		return
	}

	args := pflag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}
	for _, arg := range args {
		if arg == "-" {
			arg = "/dev/stdin"
		}

		if err := printHashes(arg, hashes); err != nil {
			panic(err)
		}
	}

}
