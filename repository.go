package main

import "path/filepath"

type Repository struct {
	dotDir FileSystemObject
}

func (repo *Repository) LoadOrInit(directory string) error {
	var err error
	repo.dotDir, err = FileSystemObject(filepath.Join(directory, ".git")).ClimbUpwardsToFind(".git")
	if nil != err {
		println("Must init")
	}
	return nil
}
