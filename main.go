package main

import (
	"flag"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/twendt/release-flow-version/pkg/buildagent"
	"github.com/twendt/release-flow-version/pkg/config"
	"github.com/twendt/release-flow-version/pkg/repository"
	"github.com/twendt/release-flow-version/pkg/usecases"
)

func main() {
	var (
		repoPath            string
		final, onBuildAgent bool
	)

	flag.StringVar(&repoPath, "p", ".", "Path to git repo")
	flag.BoolVar(&final, "final", false, "Print final version (major.minor.patch)")
	flag.BoolVar(&onBuildAgent, "buildserver", false, "Running on build agent")

	flag.Parse()

	r, err := git.PlainOpen(repoPath)
	if err != nil {
		panic(err)
	}

	err = config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to load config: %s", err))
	}

	cfg := config.GetConfig()

	repo := repository.NewGitRepo(r, cfg)

	currentBranch := &repository.Branch{}
	if !onBuildAgent {
		currentBranch, err = repo.CurrentBranch()
		if err != nil {
			panic(fmt.Errorf("Failed to get current branch: %s", err))
		}
	} else {
		ba, err := buildagent.Resolve()
		if err != nil {
			panic(fmt.Errorf("No build agent found: %s", err))
		}
		branchName, err := ba.BranchName()
		if err != nil {
			panic(fmt.Errorf("Could not get branch name fro build agent: %s", err))
		}

		branchName = repository.TrimRefPrefix(branchName, cfg.RemoteName)
		currentBranch, err = repo.FindBranchByName(branchName)
		if err != nil {
			panic(fmt.Errorf("Branch %s not found: %s", branchName, err))
		}
	}

	var version *semver.Version
	switch currentBranch.BranchType() {
	case repository.BranchTypeRelease:
		version, err = usecases.ReleaseBranchVersion(repo, currentBranch)
	case repository.BranchTypeMain:
		version, err = usecases.MainBranchVersion(repo, cfg)
	case repository.BranchTypeFeature:
		version, err = usecases.FeatureOrFixBranchVersion(repo, currentBranch)
	case repository.BranchTypeHotfix:
		version, err = usecases.FeatureOrFixBranchVersion(repo, currentBranch)
	default:
		panic(fmt.Errorf("Unsupported branch %s", currentBranch.RawName()))
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
