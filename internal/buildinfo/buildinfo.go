// Package buildinfo reports the version the binary was built from.
package buildinfo

import "runtime/debug"

// DevelVersion is reported when no version was stamped and the module has none.
const DevelVersion = "dev"

// version is stamped at release time with
// -ldflags "-X github.com/eshiji/nana/internal/buildinfo.version=vX.Y.Z".
var version = ""

// buildInfoReader matches debug.ReadBuildInfo so tests can inject build metadata.
type buildInfoReader func() (*debug.BuildInfo, bool)

// Version returns the stamped version, then the module version (go install), then DevelVersion.
func Version() string {
	return resolve(version, debug.ReadBuildInfo)
}

func resolve(stamped string, read buildInfoReader) string {
	if stamped != "" {
		return stamped
	}
	if info, ok := read(); ok && isRealVersion(info.Main.Version) {
		return info.Main.Version
	}
	return DevelVersion
}

func isRealVersion(v string) bool {
	return v != "" && v != "(devel)"
}
