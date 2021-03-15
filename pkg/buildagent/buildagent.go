package buildagent

import "fmt"

type BuildAgent interface {
	Found() bool
	BranchName() (string, error)
}

var buildAgents []BuildAgent

func Resolve() (BuildAgent, error) {
	for _, b := range buildAgents {
		if b.Found() {
			return b, nil
		}
	}

	return nil, fmt.Errorf("No buid agent found in current environment")
}
