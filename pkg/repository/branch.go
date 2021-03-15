package repository

import (
	"fmt"
	"github.com/twendt/release-flow-version/pkg/config"
	"strings"
)

type BranchType int

const (
	branchTypeUnsupported BranchType = iota
	BranchTypeMain
	BranchTypeRelease
	BranchTypeFeature
	BranchTypeHotfix
)

type Branch struct {
	name   string
	Remote string
	cfg    *config.Config
}

func NewBranch(cfg *config.Config, name string, remote string) *Branch {
	return &Branch{name: name, Remote: remote, cfg: cfg}
}

func (b *Branch) IsReleaseBranch() bool {
	remoteName := b.cfg.RemoteName
	name := strings.TrimPrefix(b.name, "refs/heads/")
	name = strings.TrimPrefix(name, fmt.Sprintf("refs/remotes/%s/", remoteName))
	releaseRegex := b.cfg.ReleaseRegex
	return releaseRegex.MatchString(name)
}

func (b *Branch) IsMainBranch() bool {
	name := TrimRefPrefix(b.ShortName(), b.cfg.RemoteName)
	return b.cfg.MainRegex.MatchString(name)
	//return b.ShortName() == b.cfg.MainBranch
}

func (b *Branch) IsFeatureBranch() bool {
	remoteName := b.cfg.RemoteName
	name := strings.TrimPrefix(b.name, "refs/heads/")
	name = strings.TrimPrefix(name, fmt.Sprintf("refs/remotes/%s/", remoteName))
	featureRegex := b.cfg.FeatureRegex
	return featureRegex.MatchString(name)
}

func (b *Branch) IsHotfixBranch() bool {
	name := TrimRefPrefix(b.name, b.cfg.RemoteName)
	return b.cfg.HotfixRegex.MatchString(name)
}

func (b *Branch) BranchType() BranchType {
	switch {
	case b.IsMainBranch():
		return BranchTypeMain
	case b.IsReleaseBranch():
		return BranchTypeRelease
	case b.IsFeatureBranch():
		return BranchTypeFeature
	case b.IsHotfixBranch():
		return BranchTypeHotfix
	}
	return branchTypeUnsupported
}

func (b *Branch) ShortName() string {
	i := strings.LastIndex(b.name, "/")

	return b.name[i+1:]
}

func (b *Branch) Name() string {
	remote := b.cfg.RemoteName
	var remotePrefix = fmt.Sprint("refs/heads/remotes/%s/", remote)
	var localPrefix = "refs/heads/"
	name := strings.TrimPrefix(b.name, remotePrefix)
	name = strings.TrimPrefix(name, localPrefix)

	return name
}

func (b *Branch) RawName() string {
	return b.name
}

//TrimRefPrefix removes the head and remote prefixes from branchName
//If remoteName is empty then origin is used
func TrimRefPrefix(branchName, remoteName string) string {
	if remoteName == "" {
		remoteName = "origin"
	}
	name := strings.TrimPrefix(branchName, "refs/heads/")
	name = strings.TrimPrefix(name, fmt.Sprintf("refs/remotes/%s/", remoteName))

	return name
}
