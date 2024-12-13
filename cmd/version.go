package cmd

// version holds the current version of the application.
// Manually update this variable for each new release.
var version = "v0.1.2"

// GetVersion returns the current version.
func GetVersion() string {
    return version
}
