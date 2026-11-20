package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteThenCheckSucceeds(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cli")
	var stderr bytes.Buffer

	if code := run([]string{dir}, &stderr); code != 0 {
		t.Fatalf("write: exit %d: %s", code, stderr.String())
	}
	if code := run([]string{"-check", dir}, &stderr); code != 0 {
		t.Fatalf("check after write: exit %d: %s", code, stderr.String())
	}
}

func TestCheckFailsOnStaleReference(t *testing.T) {
	tests := map[string]func(t *testing.T, dir string){
		"edited file":  func(t *testing.T, dir string) { appendTo(t, firstFile(t, dir), "edited\n") },
		"extra file":   func(t *testing.T, dir string) { appendTo(t, filepath.Join(dir, "extra.md"), "x") },
		"missing file": func(t *testing.T, dir string) { removeFile(t, firstFile(t, dir)) },
	}
	for name, makeStale := range tests {
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "cli")
			var stderr bytes.Buffer
			run([]string{dir}, &stderr)
			makeStale(t, dir)

			if code := run([]string{"-check", dir}, &stderr); code != 1 {
				t.Fatalf("exit %d, want 1", code)
			}
			if !bytes.Contains(stderr.Bytes(), []byte("out of date")) {
				t.Errorf("stderr %q does not explain the failure", stderr.String())
			}
		})
	}
}

func TestCheckFailsWhenDirectoryIsMissing(t *testing.T) {
	var stderr bytes.Buffer
	if code := run([]string{"-check", filepath.Join(t.TempDir(), "absent")}, &stderr); code != 1 {
		t.Fatalf("exit %d, want 1", code)
	}
}

func TestUnknownFlagIsUsageError(t *testing.T) {
	var stderr bytes.Buffer
	if code := run([]string{"-nope"}, &stderr); code != 2 {
		t.Fatalf("exit %d, want 2", code)
	}
}

func firstFile(t *testing.T, dir string) string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) == 0 {
		t.Fatalf("no generated files in %s: %v", dir, err)
	}
	return filepath.Join(dir, entries[0].Name())
}

func appendTo(t *testing.T, path, text string) {
	t.Helper()
	root, err := os.OpenRoot(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = root.Close() }()
	f, err := root.OpenFile(filepath.Base(path), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	if _, err := f.WriteString(text); err != nil {
		t.Fatal(err)
	}
}

func removeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
}
