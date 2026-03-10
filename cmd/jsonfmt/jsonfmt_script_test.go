package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	exitVal := testscript.RunMain(m, map[string]func() int{
		"jsonfmt": func() int {
			main()
			return 0
		},
	})
	os.Exit(exitVal)
}

func TestJsonfmtScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Files: []string{"../../testdata/jsonfmt.txtar"},
	})
}
