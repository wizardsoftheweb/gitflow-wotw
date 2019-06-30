package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func IsRefNameValid(ref_name string) bool {
	pattern := regexp.MustCompile(`(?m)^.*?/\..*?$|^.*?\.(lock)?$|^[^/]+$|^.*?\.\..*?$|^.*?[\000-\037\177 ~^:?*[]+.*?$|^\..*?$|^.*?/$|^.*?//.*?$|^.*?@\{.*?$|^@$|^.*?\\.*?$`)
	result := pattern.MatchString(ref_name)
	return !result
}

func ListFilesInRefsHeads(root_path string) []string {
	var ref_paths []string
	err := filepath.Walk(root_path, func(path string, info os.FileInfo, err error) error {
		info, stat_err := os.Stat(path)
		CheckIfError(stat_err)
		if info.IsDir() {
			return nil
		}
		ref_paths = append(ref_paths, path)
		return nil
	})
	CheckIfError(err)
	return ref_paths
}

func DumpFileContents(path string) {
	data, err := ioutil.ReadFile(path)
	CheckIfError(err)
	fmt.Printf("File %s: %s\n", path, string(data))
}

func ParseRevisionWithGoGit(ref_name string) string {
	return ""
}

func ParseRevisionWithGit(ref_name string) {
	// https://git-scm.com/docs/git-check-ref-format
	command := exec.Command("git", "check-ref-format", ref_name)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if nil != err {
		Warning(fmt.Sprintf("error: %s", err))
	}
}

func ListAndProcessFiles(repo_directory string) {
	head_path := filepath.Join(repo_directory, ".git", "refs", "heads")
	ref_paths := ListFilesInRefsHeads(head_path)
	for _, ref_path := range ref_paths {
		ref_name := strings.ReplaceAll(ref_path, head_path, "")
		fmt.Println(ref_name)
		ParseRevisionWithGit(ref_name)
	}
}

func MakeACommit(repo *git.Repository, directory string, random_stuff string) {
	// Pulled from the commit demo
	w, err := repo.Worktree()
	CheckIfError(err)
	Info("echo \"hello world!\" > example-git-file")
	filename := filepath.Join(directory, "example-git-file")
	err = ioutil.WriteFile(filename, []byte(random_stuff), 0644)
	CheckIfError(err)
	Info("git add example-git-file")
	_, err = w.Add("example-git-file")
	CheckIfError(err)
	Info("git status --porcelain")
	status, err := w.Status()
	CheckIfError(err)
	fmt.Println(status)
	Info("git commit -m \"example go-git commit\"")
	commit, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})
	CheckIfError(err)
	Info("git show -s")
	obj, err := repo.CommitObject(commit)
	CheckIfError(err)
	fmt.Println(obj)
}

func CreateNewBranch(repo *git.Repository, branch_name string) {
	// Pulled from the branch demo
	repo_config, err := repo.Config()
	CheckIfError(err)
	headRef, err := repo.Head()
	CheckIfError(err)
	new_ref := plumbing.NewHashReference(plumbing.NewBranchReferenceName(branch_name), headRef.Hash())
	fmt.Printf(
		"Ref name %s is a valid name: %t\n",
		new_ref.Name().String(),
		IsRefNameValid(new_ref.Name().String()),
	)
	err = repo.Storer.SetReference(new_ref)
	CheckIfError(err)
	err = repo.Storer.SetConfig(repo_config)
	CheckIfError(err)
	w, err := repo.Worktree()
	CheckIfError(err)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: new_ref.Name(),
	})
	CheckIfError(err)
	MakeACommit(repo, w.Filesystem.Root(), branch_name)
}

func RunProcess() {
	directory, _ := os.Getwd()
	r, err := git.PlainInit(directory, false)
	if git.ErrRepositoryAlreadyExists == err {
		r, err = git.PlainOpen(directory)
	}
	CheckIfError(err)
	MakeACommit(r, directory, "master")
	CreateNewBranch(r, "new-..bad\\.branch//name.")
}

func main() {
	// Sets up the repo
	RunProcess()
	directory, _ := os.Getwd()
	r, err := git.PlainOpen(directory)
	CheckIfError(err)
	repo_config, err := r.Config()
	CheckIfError(err)
	for branch_name, branches := range repo_config.Branches {
		fmt.Println(branch_name, branches)
	}

	ListAndProcessFiles(directory)

	// Or deleted from it.
	// Info("git branch -D my-branch")
	// err = r.Storer.RemoveReference(ref.Name())
	// CheckIfError(err)
}
