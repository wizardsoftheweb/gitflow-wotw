package gitflow

import (
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	GitBranchPattern        = regexp.MustCompile(`(?m)^.*?(\w.*)\s*?`)
	GitCurrentBranchPattern = regexp.MustCompile(`(?m)^\s*\*\s*?(\w.*)\s*?$`)
)

type Repository struct {
	Prefix      string
	HumanPrefix string
}

var (
	Repo = &Repository{}
)

func (r *Repository) SpecificPrefixBranches(remote bool) []string {
	var prefixedBranches []string
	for _, branch := range r.SpecificBranches(remote) {
		if strings.HasPrefix(branch, r.Prefix) {
			prefixedBranches = append(prefixedBranches, strings.Replace(branch, r.Prefix, "", 1))
		}
	}
	return prefixedBranches
}

func (r *Repository) SpecificBranches(remote bool) []string {
	var branches []string
	result := BranchNoColor(remote)
	for _, match := range GitBranchPattern.FindAllStringSubmatch(result.result, -1) {
		branches = append(branches, match[1])
	}
	return branches
}

func (r *Repository) LocalBranches() []string {
	return r.SpecificBranches(false)
}

func (r *Repository) RemoteBranches() []string {
	return r.SpecificBranches(true)
}

func (r *Repository) CurrentBranch() string {
	results := GitCurrentBranchPattern.FindStringSubmatch(BranchNoColor(false).String())
	if 2 < len(results) {
		logrus.Fatal(ErrCannotDetermineCurrentBranch)
	}
	return results[1]
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

func (r *Repository) PickGoodSuggestion(branchName string) string {
	if "master" == branchName {
		return r.PickGoodMasterSuggestion()
	} else {
		return r.PickGoodDevSuggestion()
	}
}

func (r *Repository) PickGoodMasterSuggestion() string {
	logrus.Trace("PickGoodMasterSuggestion")
	for _, suggestion := range DefaultMasterSuggestions {
		if r.HasLocalBranch(suggestion) {
			return suggestion
		}
	}
	return DefaultBranchMaster.Value
}
func (r *Repository) PickGoodDevSuggestion() string {
	logrus.Trace("PickGoodDevSuggestion")
	newMaster := GitConfig.Get(MasterBranchKey)
	for _, suggestion := range DefaultDevSuggestions {
		if suggestion != newMaster && r.HasLocalBranch(suggestion) {
			return suggestion
		}
	}
	return ""
}
