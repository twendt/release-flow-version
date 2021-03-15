package buildagent

import (
	"fmt"
	"os"
)

const (
	branchVariable = "BUILD_SOURCEBRANCH"
	agentVariable  = "TF_BUILD"
)

type AzureDevOpsAgent struct {
}

func (a *AzureDevOpsAgent) Found() bool {
	v := os.Getenv(agentVariable)
	if v == "" {
		return false
	}

	return true
}

func (a *AzureDevOpsAgent) BranchName() (string, error) {
	b := os.Getenv(branchVariable)
	if b == "" {
		return "", fmt.Errorf("Variable %s not found", branchVariable)
	}

	return b, nil
}

func init() {
	buildAgents = append(buildAgents, &AzureDevOpsAgent{})
}
