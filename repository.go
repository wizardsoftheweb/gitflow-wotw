package main

import "path/filepath"

type Repository struct {
	dotDir FileSystemObject
}

func (repo *Repository) LoadOrInit(directory string) error {
	repo.dotDir = FileSystemObject(filepath.Join(directory, ".git"))
	if repo.dotDir.exists() {
		println("rad")
	} else {
		println("git init")
	}
	return nil
}
