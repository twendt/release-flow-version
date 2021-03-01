package main

import (
	"fmt"
	"github.com/Masterminds/semver"
	"sort"
)

type release struct {
	name   string
	branch *branch
}

func newRelease(name string, branch *branch) *release {
	return &release{name: name, branch: branch}
}

func latestReleaseFromList(releases []*release) (*release, error) {
	vs := make([]*semver.Version, len(releases))
	for i, r := range releases {
		v, err := semver.NewVersion(r.name)
		if err != nil {
			return nil, fmt.Errorf("Error parsing version: %s", err)
		}

		vs[i] = v
	}

	sort.Sort(semver.Collection(vs))

	result := &release{}
	latest := vs[len(vs)-1].String()
	for _, r := range releases {
		if r.name == latest {
			result = r
		}
	}

	return result, nil
}
