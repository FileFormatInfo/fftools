package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	exitVal := testscript.RunMain(m, map[string]func() int{
		"luhncheck": func() int {
			main()
			return 0
		},
	})
	os.Exit(exitVal)
}

func TestLuhncheckScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Files: []string{"../../testdata/luhncheck.txtar"},
	})
}
