package version

import (
	"strconv"

	"github.com/coreos/go-semver/semver"
)

var (
	major    string
	minor    string
	patch    string
	pre      string
	metadata string

	Version = semver.Version{
		Major:      parseVersionNumber(major),
		Minor:      parseVersionNumber(minor),
		Patch:      parseVersionNumber(patch),
		PreRelease: semver.PreRelease(pre),
		Metadata:   metadata,
	}
)

func parseVersionNumber(versionNum string) int64 {
	if versionNum == "" {
		return 0
	}
	i, err := strconv.ParseInt(versionNum, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}
