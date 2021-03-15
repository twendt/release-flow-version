package repository

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/twendt/release-flow-version/pkg/config"
	"sort"
	"strings"
)

const (
	localBranchPrefix  = "refs/heads/"
	remoteBranchPrefix = "refs/remotes/"
)

type GitRepo struct {
	r          *git.Repository
	mainBranch *Branch
	cfg        *config.Config
}

func NewGitRepo(r *git.Repository, cfg *config.Config) *GitRepo {
	repo := &GitRepo{r: r, cfg: cfg}
	return repo
}

func (g *GitRepo) CurrentBranch() (*Branch, error) {
	head, err := g.r.Head()
	if err != nil {
		return nil, fmt.Errorf("Failed to get head: %s", err)
	}

	return NewBranch(g.cfg, string(head.Name()), ""), nil
}

func (g *GitRepo) CommitCountSinceRelease(release *Release) (int, error) {
	//mainBranchName := fmt.Sprintf("refs/heads/%s", g.cfg.MainBranch)
	mainBranch, err := g.MainBranch()
	if err != nil {
		return 0, fmt.Errorf("Failed to get main branch: %s", err)
	}

	//mainBranchName := fmt.Sprintf("refs/remotes/%s/%s", g.cfg.RemoteName, g.cfg.MainBranch)
	releaseBranchName := release.Branch.RawName()
	baseCommit, err := g.MergeBase(mainBranch.RawName(), releaseBranchName)
	if err != nil {
		return 0, fmt.Errorf("Failed to get merge base commit for %s and %s: %s", mainBranch.RawName(), release.Branch.RawName(), err)
	}

	log, err := g.r.Log(&git.LogOptions{})
	if err != nil {
		return 0, fmt.Errorf("Failed to get git log: %s", err)
	}

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

	if commits == nil || len(commits) < 2 {
		return nil, fmt.Errorf("Missing commits to find merge base")
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

func (g *GitRepo) Branches() ([]*Branch, error) {
	result := []*Branch{}

	references, err := g.r.References()
	if err != nil {
		return nil, fmt.Errorf("Failed to get references: %s", err)
	}
	err = references.ForEach(func(reference *plumbing.Reference) error {
		name := string(reference.Name())
		if strings.HasPrefix(name, localBranchPrefix) {
			b := NewBranch(g.cfg, name, "")
			result = append(result, b)
			return nil
		}

		remoteName := g.cfg.RemoteName
		if remoteName == "" {
			remoteName = "origin"
		}

		prefix := fmt.Sprintf("%s%s/", remoteBranchPrefix, remoteName)
		if strings.HasPrefix(name, prefix) {
			b := NewBranch(g.cfg, name, remoteName)
			result = append(result, b)
			return nil
		}

		return nil
	})
	references.Close()

	return result, nil
}

func (g *GitRepo) MainBranch() (*Branch, error) {
	if g.mainBranch != nil {
		return g.mainBranch, nil
	}

	branches, err := g.Branches()
	if err != nil {
		return nil, fmt.Errorf("Failed to get branches for repo: %s", err)
	}

	for _, b := range branches {
		return b, nil
	}

	return nil, fmt.Errorf("Main banch not found")
}

func (g *GitRepo) Releases() ([]*Release, error) {
	branches, err := g.Branches()
	if err != nil {
		return nil, fmt.Errorf("Failed to get branches for repo: %s", err)
	}

	releases := []*Release{}

	for _, b := range branches {
		if !b.IsReleaseBranch() {
			continue
		}

		releases = append(releases, NewRelease(b.ShortName(), b))
	}

	return releases, nil
}

func (g *GitRepo) LatestRelease() (*Release, error) {
	releases, err := g.Releases()
	if err != nil {
		return nil, fmt.Errorf("Failed to get releases: %s", err)
	}

	if len(releases) == 0 {
		return nil, nil
	}

	latestRelease, err := LatestReleaseFromList(releases)
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

func LatestReleaseFromList(releases []*Release) (*Release, error) {
	vs := make([]*semver.Version, len(releases))
	for i, r := range releases {
		v, err := semver.NewVersion(r.Name)
		if err != nil {
			return nil, fmt.Errorf("Error parsing version: %s", err)
		}

		vs[i] = v
	}

	sort.Sort(semver.Collection(vs))

	result := &Release{}
	latest := vs[len(vs)-1].String()
	for _, r := range releases {
		if r.Name == latest {
			result = r
		}
	}

	return result, nil
}

//FindBranchByName searches all branches to find a matching name
//The match can begin with refs/heads/ or refs/remotes/name/
func (g *GitRepo) FindBranchByName(name string) (*Branch, error) {
	branches, err := g.Branches()
	if err != nil {
		return nil, fmt.Errorf("Failed to get branches in repo: %s", err)
	}

	for _, b := range branches {
		if b.Name() == name {
			return b, nil
		}
	}

	return nil, fmt.Errorf("No local or remote branch %s found", name)
}
