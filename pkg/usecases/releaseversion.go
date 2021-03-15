package usecases

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/twendt/release-flow-version/pkg/repository"
)

func ReleaseBranchVersion(r *repository.GitRepo, b *repository.Branch) (*semver.Version, error) {
	v := &semver.Version{}

	release := repository.NewRelease(b.ShortName(), b)
	count, err := r.CommitCountSinceRelease(release)
	if err != nil {
		return v, fmt.Errorf("Failed to get commit count since release: %s", err)
	}

	v, err = semver.NewVersion(release.Name)
	if err != nil {
		return v, fmt.Errorf("Release branch %s has wrong name: %s", b.RawName(), err)
	}

	for i := 0; i < count; i++ {
		*v = v.IncPatch()
	}

	return v, nil
}
