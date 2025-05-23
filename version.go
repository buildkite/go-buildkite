package buildkite

import (
	"runtime/debug"
	"strings"
)

// Version returns the library version number based on Go module information
var Version = func() string {
	// Try to get version from build info
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev" // No build info available
	}
	
	// For the main module
	if info.Main.Path == "github.com/buildkite/go-buildkite/v4" && info.Main.Version != "(devel)" {
		// Remove the 'v' prefix if present
		return strings.TrimPrefix(info.Main.Version, "v")
	}
	
	// Try to find this module in the dependency list
	for _, dep := range info.Deps {
		if dep.Path == "github.com/buildkite/go-buildkite/v4" && dep.Version != "(devel)" {
			return strings.TrimPrefix(dep.Version, "v")
		}
	}
	
	// If we're in development mode, try to get VCS info
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			// Return shortened git commit hash
			if len(setting.Value) >= 7 {
				return "dev+" + setting.Value[:7]
			}
			return "dev+" + setting.Value
		}
	}
	
	// Last resort fallback
	return "dev"
}()
