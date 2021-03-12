package main

import (
	"fmt"
	"github.com/Masterminds/semver"
)

func releaseBranchVersion(r *GitRepo, b *branch) (*semver.Version, error) {
	v := &semver.Version{}

	release := newRelease(b.shortName(), b)
	count, err := r.CommitCountSinceRelease(release)
	if err != nil {
		return v, fmt.Errorf("Failed to get commit count since release: %s", err)
	}

	v, err = semver.NewVersion(release.name)
	if err != nil {
		return v, fmt.Errorf("Release branch %s has wrong name: %s", b.name, err)
	}

	for i := 0; i < count; i++ {
		*v = v.IncPatch()
	}

	return v, nil
}
