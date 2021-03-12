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
	cfg    *Config
}

func newBranch(cfg *Config, name string, remote string) *branch {
	return &branch{name: name, Remote: remote, cfg: cfg}
}

func (b *branch) isReleaseBranch() bool {
	remoteName := b.cfg.RemoteName
	name := strings.TrimPrefix(b.name, "refs/heads/")
	name = strings.TrimPrefix(name, fmt.Sprintf("refs/remotes/%s/", remoteName))
	releaseRegex := b.cfg.ReleaseRegex
	return releaseRegex.MatchString(name)
}

func (b *branch) isMainBranch() bool {
	return b.shortName() == b.cfg.MainBranch
}

func (b *branch) isFeatureBranch() bool {
	remoteName := b.cfg.RemoteName
	name := strings.TrimPrefix(b.name, "refs/heads/")
	name = strings.TrimPrefix(name, fmt.Sprintf("refs/remotes/%s/", remoteName))
	featureRegex := b.cfg.FeatureRegex
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
	remote := b.cfg.RemoteName
	var remotePrefix = fmt.Sprint("refs/heads/remotes/%s/", remote)
	var localPrefix = "refs/heads/"
	name := strings.TrimPrefix(b.name, remotePrefix)
	name = strings.TrimPrefix(name, localPrefix)

	return name
}
