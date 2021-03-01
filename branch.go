package main

import (
	"fmt"
	"strings"
)

type branchType int

const (
	branchTypeUnsupported branchType = iota
	branchTypeMain
	branchTypeRelease
	branchTypeFeature
)

type branch struct {
	name   string
	Remote string
}

func newBranch(name string, remote string) *branch {
	return &branch{name: name, Remote: remote}
}

func (b *branch) isReleaseBranch() bool {
	releaseRegex := GetConfig().ReleaseRegex
	return releaseRegex.MatchString(b.name)
}

func (b *branch) isMainBranch() bool {
	return b.shortName() == GetConfig().MainBranch
}

func (b *branch) isFeatureBranch() bool {
	featureRegex := GetConfig().FeatureRegex
	return featureRegex.MatchString(b.name)
}

func (b *branch) branchType() branchType {
	switch {
	case b.isMainBranch():
		return branchTypeMain
	case b.isReleaseBranch():
		return branchTypeRelease
	case b.isFeatureBranch():
		return branchTypeFeature
	}
	return branchTypeUnsupported
}

func (b *branch) shortName() string {
	i := strings.LastIndex(b.name, "/")

	return b.name[i+1:]
}

func (b *branch) Name() string {
	remote := GetConfig().RemoteName
	var remotePrefix = fmt.Sprint("refs/heads/remotes/%s/", remote)
	var localPrefix = "refs/heads/"
	name := strings.TrimPrefix(b.name, remotePrefix)
	name = strings.TrimPrefix(name, localPrefix)

	return name
}
