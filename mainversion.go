package main

import (
	"fmt"
)

func mainBranchVersion(r *GitRepo) (string, error) {
	latestRelease, err := r.LatestRelease()
	if err != nil {
		return "", fmt.Errorf("Failed to get latest release branch name: %s", err)
	}

	if latestRelease == nil {
		counter, err := r.CommitCountCurrentBranch()
		if err != nil {
			return "", fmt.Errorf("Failed to get commit count for current branch: %s", err)
		}
		return fmt.Sprintf("0.1.0-beta.%d", counter), nil
	}

	counter, err := r.CommitCountSinceRelease(latestRelease)
	if err != nil {
		return "", fmt.Errorf("Failed to get commit count since last release: %s", err)
	}

	return fmt.Sprintf("%s-beta.%d", latestRelease.name, counter), nil

}
