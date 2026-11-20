// Command gendocs writes the CLI reference to a directory, or with -check verifies that
// the directory is current. The reference is generated from the command tree itself.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra/doc"

	"github.com/eshiji/nana/internal/cli"
)

const defaultDir = "docs/cli"

var errStale = errors.New("CLI reference is out of date: run `make docs`")

func main() {
	os.Exit(run(os.Args[1:], os.Stderr))
}

func run(args []string, stderr io.Writer) int {
	flags := flag.NewFlagSet("gendocs", flag.ContinueOnError)
	flags.SetOutput(stderr)
	check := flags.Bool("check", false, "fail if the reference differs instead of writing it")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	dir := defaultDir
	if flags.NArg() > 0 {
		dir = flags.Arg(0)
	}
	if err := generateInto(dir, *check); err != nil {
		_, _ = fmt.Fprintln(stderr, "gendocs:", err)
		return 1
	}
	return 0
}

func generateInto(dir string, check bool) error {
	fresh, err := os.MkdirTemp("", "gendocs-")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(fresh) }()
	if err := generate(fresh); err != nil {
		return err
	}
	if check {
		return compare(fresh, dir)
	}
	return replace(fresh, dir)
}

func generate(dir string) error {
	root := cli.NewRootCommand(io.Discard, io.Discard)
	root.DisableAutoGenTag = true // keeps the output deterministic
	if err := doc.GenMarkdownTree(root, dir); err != nil {
		return fmt.Errorf("generating markdown: %w", err)
	}
	return nil
}

func compare(fresh, current string) error {
	want, err := readTree(fresh)
	if err != nil {
		return err
	}
	got, err := readTree(current)
	if err != nil {
		return err
	}
	if !equalTrees(want, got) {
		return errStale
	}
	return nil
}

func replace(fresh, dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("removing %s: %w", dir, err)
	}
	if err := os.MkdirAll(filepath.Dir(dir), 0o750); err != nil {
		return fmt.Errorf("creating parent of %s: %w", dir, err)
	}
	if err := os.CopyFS(dir, os.DirFS(fresh)); err != nil {
		return fmt.Errorf("writing %s: %w", dir, err)
	}
	return nil
}

// readTree reads every file in dir through an os.Root, so no entry can escape dir.
func readTree(dir string) (map[string][]byte, error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", dir, err)
	}
	defer func() { _ = root.Close() }()
	entries, err := fs.ReadDir(root.FS(), ".")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}
	tree := make(map[string][]byte, len(entries))
	for _, entry := range entries {
		content, err := root.ReadFile(entry.Name())
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", entry.Name(), err)
		}
		tree[entry.Name()] = content
	}
	return tree, nil
}

func equalTrees(a, b map[string][]byte) bool {
	if len(a) != len(b) {
		return false
	}
	for name, content := range a {
		other, ok := b[name]
		if !ok || !bytes.Equal(content, other) {
			return false
		}
	}
	return true
}
