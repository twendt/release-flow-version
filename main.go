package main

import (
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
)

func main() {
	var repoPath string

	flag.StringVar(&repoPath, "p", ".", "Path to git repo")
	flag.Parse()

	r, err := git.PlainOpen(repoPath)
	if err != nil {
		panic(err)
	}

	err = LoadConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to load config: %s", err))
	}

	repo := NewGitRepo(r, GetConfig())

	currentBranch, err := repo.CurrentBranch()
	if err != nil {
		panic(fmt.Errorf("Failed to get current branch: %s", err))
	}

	var version string
	switch currentBranch.branchType() {
	case branchTypeRelease:
		version, err = releaseBranchVersion(repo, currentBranch)
	case branchTypeMain:
		version, err = mainBranchVersion(repo)
	case branchTypeFeature:
		version, err = featureBranchVersion(repo, currentBranch)
	default:
		panic(fmt.Errorf("Unsupported branch %s", currentBranch.name))
	}

	if err != nil {
		panic(err)
	}

	fmt.Println(version)

}
