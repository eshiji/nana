package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	tests := map[string]struct {
		args       []string
		wantCode   int
		wantStdout bool
		wantStderr string
	}{
		"version prints a version":    {[]string{"version"}, ExitOK, true, ""},
		"version rejects arguments":   {[]string{"version", "extra"}, ExitError, false, "unknown command"},
		"unknown command is an error": {[]string{"does-not-exist"}, ExitError, false, "unknown command"},
		"no arguments shows help":     {nil, ExitOK, true, ""},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer

			code := Execute(tc.args, &stdout, &stderr)

			if code != tc.wantCode {
				t.Errorf("exit code = %d, want %d (stderr: %s)", code, tc.wantCode, stderr.String())
			}
			if got := stdout.Len() > 0; got != tc.wantStdout {
				t.Errorf("stdout written = %v, want %v", got, tc.wantStdout)
			}
			if !strings.Contains(stderr.String(), tc.wantStderr) {
				t.Errorf("stderr = %q, want it to contain %q", stderr.String(), tc.wantStderr)
			}
		})
	}
}
