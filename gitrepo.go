package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"strings"
)

const (
	localBranchPrefix  = "refs/heads/"
	remoteBranchPrefix = "refs/remotes/"
)

type GitRepo struct {
	r   *git.Repository
	cfg *Config
}

func NewGitRepo(r *git.Repository, cfg *Config) *GitRepo {
	return &GitRepo{r: r, cfg: cfg}
}

func (g *GitRepo) CurrentBranch() (*branch, error) {
	head, err := g.r.Head()
	if err != nil {
		return nil, fmt.Errorf("Failed to get head: %s", err)
	}

	return newBranch(string(head.Name()), ""), nil
}

func (g *GitRepo) CommitCountSinceRelease(release *release) (int, error) {
	mainBranchName := fmt.Sprintf("refs/heads/%s", g.cfg.MainBranch)
	releaseBranchName := release.branch.name
	baseCommit, err := g.MergeBase(mainBranchName, releaseBranchName)
	if err != nil {
		return 0, fmt.Errorf("Failed to get merge base commit for %s and %s: %s", g.cfg.MainBranch, release.branch.name, err)
	}

	log, err := g.r.Log(&git.LogOptions{})

	counter := 0
	for {
		c, err := log.Next()
		if err != nil {
			return 0, fmt.Errorf("Failed to traverse commits")
		}

		if c.Hash == baseCommit.Hash {
			break
		}

		counter++
	}
	return counter, nil
}

func (g *GitRepo) MergeBase(b1, b2 string) (*object.Commit, error) {
	var hashes []*plumbing.Hash
	for _, rev := range []string{b1, b2} {
		hash, err := g.r.ResolveRevision(plumbing.Revision(rev))
		if err != nil {
			return nil, fmt.Errorf("could not parse revision '%s': %s", rev, err)
		}
		hashes = append(hashes, hash)
	}

	// Get the commits identified by the passed hashes
	var commits []*object.Commit
	for _, hash := range hashes {
		commit, err := g.r.CommitObject(*hash)
		if err != nil {
			return nil, fmt.Errorf("could not find commit '%s': %s", hash.String(), err)
		}
		commits = append(commits, commit)
	}

	res, err := commits[0].MergeBase(commits[1])
	if err != nil {
		return nil, fmt.Errorf("could not traverse the repository history: %s", err)
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("Could not find merge base for %s and %s", b1, b2)
	}
	return res[0], nil
}

func (g *GitRepo) Branches() ([]*branch, error) {
	result := []*branch{}

	references, err := g.r.References()
	if err != nil {
		return nil, fmt.Errorf("Failed to get references: %s", err)
	}
	err = references.ForEach(func(reference *plumbing.Reference) error {
		name := string(reference.Name())
		if strings.HasPrefix(name, localBranchPrefix) {
			b := newBranch(name, "")
			result = append(result, b)
			return nil
		}

		remoteName := cfg.RemoteName
		if remoteName == "" {
			remoteName = "origin"
		}

		prefix := fmt.Sprintf("%s%s/", remoteBranchPrefix, remoteName)
		if strings.HasPrefix(name, prefix) {
			b := newBranch(name, remoteName)
			result = append(result, b)
			return nil
		}

		return nil
	})
	references.Close()

	return result, nil
}

func (g *GitRepo) Releases() ([]*release, error) {
	branches, err := g.Branches()
	if err != nil {
		return nil, fmt.Errorf("Failed to get branches for repo: %s", err)
	}

	releases := []*release{}

	for _, b := range branches {
		if !b.isReleaseBranch() {
			continue
		}

		releases = append(releases, newRelease(b.shortName(), b))
	}

	return releases, nil
}

func (g *GitRepo) LatestRelease() (*release, error) {
	releases, err := g.Releases()
	if err != nil {
		return nil, fmt.Errorf("Failed to get releases: %s", err)
	}

	if len(releases) == 0 {
		return nil, nil
	}

	latestRelease, err := latestReleaseFromList(releases)
	if err != nil {
		return nil, fmt.Errorf("Failed to get latest release from %v: %s", releases, err)
	}

	return latestRelease, nil
}

func (g *GitRepo) CommitCountCurrentBranch() (int, error) {
	log, err := g.r.Log(&git.LogOptions{})
	if err != nil {
		return 0, fmt.Errorf("Failed to get log: %s", err)
	}

	counter := 0
	err = log.ForEach(func(commit *object.Commit) error {
		counter++
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("Failed to count commits: %s", err)
	}

	return counter, nil
}
