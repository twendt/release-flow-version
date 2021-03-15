package repository

type Release struct {
	Name   string
	Branch *Branch
}

func NewRelease(name string, branch *Branch) *Release {
	return &Release{Name: name, Branch: branch}
}
