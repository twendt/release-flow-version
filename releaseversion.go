package main

import (
	"fmt"
	"github.com/Masterminds/semver"
)

func releaseBranchVersion(r *GitRepo, b *branch) (string, error) {
	release := newRelease(b.shortName(), b)
	count, err := r.CommitCountSinceRelease(release)
	if err != nil {
		return "", fmt.Errorf("Failed to get commit count since release: %s", err)
	}

	base, err := semver.NewVersion(release.name)
	if err != nil {
		return "", fmt.Errorf("Release branch %s has wrong name: %s", b.name, err)
	}

	for i := 0; i < count; i++ {
		*base = base.IncPatch()
	}

	return base.String(), nil
}
