// Command nana is the entry point; all behaviour lives in internal/.
package main

import (
	"os"

	"github.com/eshiji/nana/internal/cli"
)

func main() {
	os.Exit(cli.Execute(os.Args[1:], os.Stdout, os.Stderr))
}
