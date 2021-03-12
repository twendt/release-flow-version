package main

import (
	"fmt"
	"github.com/Masterminds/semver"
)

func featureBranchVersion(r *GitRepo, currentBranch *branch) (*semver.Version, error) {
	v := &semver.Version{}

	latestRelease, err := r.LatestRelease()
	if err != nil {
		return v, fmt.Errorf("Failed to get latest release branch name: %s", err)
	}

	counter, err := r.CommitCountSinceRelease(latestRelease)
	if err != nil {
		return v, fmt.Errorf("Failed to get commit count since last release: %s", err)
	}

	//latestRelease := strings.TrimPrefix(latestReleaseBranch, "release/")

	v, err = semver.NewVersion(latestRelease.name)
	if err != nil {
		return v, fmt.Errorf("Release name is not a version number: %s", err)
	}

	*v = v.IncMinor()
	*v, err = v.SetPrerelease(prereleaseStr(currentBranch.shortName(), counter))
	if err != nil {
		return v, fmt.Errorf("Failed to set prerelease string: %s", err)
	}

	return v, nil
	//return fmt.Sprintf("%s.%s.%d", latestRelease.name, currentBranch.shortName(), counter), nil

}

func prereleaseStr(tag string, counter int) string {
	return fmt.Sprintf("%s.%d", tag, counter)
}
