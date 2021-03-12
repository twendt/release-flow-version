package main

import (
	"flag"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
)

func main() {
	var (
		repoPath, branchName, mainBranchName string
		final                                bool
	)

	flag.StringVar(&repoPath, "p", ".", "Path to git repo")
	flag.BoolVar(&final, "final", false, "Print final version (major.minor.patch)")
	flag.StringVar(&branchName, "branch", "", "Name of branch, if not specified then the current checked out branch is used")
	flag.StringVar(&mainBranchName, "main", "", "Name of the main branch, defaults to main")

	flag.Parse()

	r, err := git.PlainOpen(repoPath)
	if err != nil {
		panic(err)
	}

	err = LoadConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to load config: %s", err))
	}

	cfg := GetConfig()
	if mainBranchName != "" {
		cfg.MainBranch = mainBranchName
	}

	repo := NewGitRepo(r, cfg)

	currentBranch := &branch{}
	if branchName == "" {
		currentBranch, err = repo.CurrentBranch()
		if err != nil {
			panic(fmt.Errorf("Failed to get current branch: %s", err))
		}
	} else {
		currentBranch = newBranch(cfg, branchName, "")
	}

	var version *semver.Version
	switch currentBranch.branchType() {
	case branchTypeRelease:
		version, err = releaseBranchVersion(repo, currentBranch)
	case branchTypeMain:
		version, err = mainBranchVersion(repo, cfg)
	case branchTypeFeature:
		version, err = featureBranchVersion(repo, currentBranch)
	default:
		panic(fmt.Errorf("Unsupported branch %s", currentBranch.name))
	}

	if err != nil {
		panic(err)
	}

	if final {
		fmt.Printf("%d.%d.%d\n", version.Major(), version.Minor(), version.Patch())
		return
	}

	fmt.Println(version.String())
	return
}
