package main

import (
	"fmt"
	"github.com/Masterminds/semver"
)

func mainBranchVersion(r *GitRepo, cfg *Config) (*semver.Version, error) {
	v, err := semver.NewVersion(cfg.DefaultVersion)
	if err != nil {
		return v, fmt.Errorf("Default version string is not a proper version number: %s", err)
	}

	latestRelease, err := r.LatestRelease()
	if err != nil {
		return v, fmt.Errorf("Failed to get latest release branch name: %s", err)
	}

	if latestRelease == nil {
		counter, err := r.CommitCountCurrentBranch()
		if err != nil {
			return v, fmt.Errorf("Failed to get commit count for current branch: %s", err)
		}

		return versionWithPrerelease(v, prereleaseStr(cfg.MainTag, counter))
	}

	counter, err := r.CommitCountSinceRelease(latestRelease)
	if err != nil {
		return v, fmt.Errorf("Failed to get commit count since last release: %s", err)
	}

	v, err = semver.NewVersion(latestRelease.name)
	if err != nil {
		return v, fmt.Errorf("Latest release is not semver: %s", err)
	}

	*v = v.IncMinor()

	return versionWithPrerelease(v, prereleaseStr(cfg.MainTag, counter))
}

func versionWithPrerelease(v *semver.Version, pre string) (*semver.Version, error) {
	var err error
	*v, err = v.SetPrerelease(pre)
	if err != nil {
		return v, fmt.Errorf("Invalid prerelease string \"%s\": %s", pre, err)
	}

	return v, nil
}
