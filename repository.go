package main

import "strings"

type Repository struct{}

var (
	Repo = &Repository{}
)

func (r *Repository) SpecificBranches(remote bool) []string {
	branches := []string{}
	result := BranchNoColor(remote)
	for _, branch := range strings.Split(result.result, `\n`) {
		branches = append(
			branches,
			strings.TrimSpace(strings.TrimPrefix(branch, `*`)),
		)
	}
	return branches
}

func (r *Repository) LocalBranches() []string {
	return r.SpecificBranches(false)
}

func (r *Repository) RemoteBranches() []string {
	return r.SpecificBranches(true)
}

func (r *Repository) HasLocalBranch(needle string) bool {
	for _, branch := range r.LocalBranches() {
		if needle == branch {
			return true
		}
	}
	return false
}

func (r *Repository) HasRemoteBranch(needle string) bool {
	for _, branch := range r.RemoteBranches() {
		if needle == branch {
			return true
		}
	}
	return false
}

func (r *Repository) PickGoodMasterSuggestion() string {
	for _, suggestion := range DefaultMasterSuggestions {
		if r.HasLocalBranch(suggestion) {
			return suggestion
		}
	}
	return DefaultBranchMaster.Value
}
func (r *Repository) PickGoodDevSuggestion() string {
	for _, suggestion := range DefaultDevSuggestions {
		if r.HasLocalBranch(suggestion) {
			return suggestion
		}
	}
	return DefaultBranchDevelopment.Value
}
