package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	exitVal := testscript.RunMain(m, nil)

	os.Exit(exitVal)
}

func TestUrly(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "../../testdata",
	})
}
