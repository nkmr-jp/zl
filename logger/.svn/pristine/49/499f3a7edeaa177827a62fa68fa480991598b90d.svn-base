package logger

import (
	"log"
	"os/exec"
	"strings"
)

// Deprecated: it will be removed in the next major version upgrade.
type VersionType int

const (
	// Deprecated: it will be removed in the next major version upgrade.
	VersionTypeRevision VersionType = iota
	// Deprecated: it will be removed in the next major version upgrade.
	VersionTypeTag
)

// GetVersionFromGit use the git revision or tag as a version.
// When using tag, recommend semantic versioning.
// See https://semver.org/
//
// Deprecated: it will be removed in the next major version upgrade.
func GetVersion(versionType VersionType) *string {
	var out []byte
	var err error

	switch versionType {
	case VersionTypeRevision:
		out, err = exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	case VersionTypeTag:
		out, err = exec.Command("git", "tag").Output()
	}
	if err != nil {
		log.Print(err)
		return nil
	}

	ret := strings.TrimRight(string(out), "\n")
	return &ret
}
