package buildinfo

import (
	"runtime/debug"
	"testing"
)

func readerWith(mainVersion string, ok bool) buildInfoReader {
	return func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{Main: debug.Module{Version: mainVersion}}, ok
	}
}

func TestResolve(t *testing.T) {
	tests := map[string]struct {
		stamped string
		read    buildInfoReader
		want    string
	}{
		"stamped version wins":        {"v1.2.3", readerWith("v0.9.0", true), "v1.2.3"},
		"module version from install": {"", readerWith("v0.9.0", true), "v0.9.0"},
		"devel build falls back":      {"", readerWith("(devel)", true), DevelVersion},
		"empty module version":        {"", readerWith("", true), DevelVersion},
		"build info unavailable":      {"", readerWith("v0.9.0", false), DevelVersion},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := resolve(tc.stamped, tc.read); got != tc.want {
				t.Errorf("resolve() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestVersionIsNeverEmpty(t *testing.T) {
	if Version() == "" {
		t.Fatal("Version() returned an empty string")
	}
}
