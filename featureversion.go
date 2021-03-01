package main

import (
	"fmt"
)

func featureBranchVersion(r *GitRepo, currentBranch *branch) (string, error) {
	latestRelease, err := r.LatestRelease()
	if err != nil {
		return "", fmt.Errorf("Failed to get latest release branch name: %s", err)
	}

	counter, err := r.CommitCountSinceRelease(latestRelease)
	if err != nil {
		return "", fmt.Errorf("Failed to get commit count since last release: %s", err)
	}

	//latestRelease := strings.TrimPrefix(latestReleaseBranch, "release/")

	return fmt.Sprintf("%s.%s.%d", latestRelease.name, currentBranch.shortName(), counter), nil
}
