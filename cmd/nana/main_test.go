package main

import (
	"flag"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

var update = flag.Bool("update", false, "rewrite testscript golden files (make golden)")

// TestMain lets testscript run the real binary in-process as "nana".
func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"nana": main,
	})
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{Dir: "testdata/script", UpdateScripts: *update})
}
