package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

type FileSystemObject string

var (
	ErrRootReached = errors.New("Root directory reached without finding file")
)

func (path FileSystemObject) String() string {
	return string(path)
}

func (path FileSystemObject) Parent() FileSystemObject {
	return FileSystemObject(filepath.Dir(path.String()))
}

func (path FileSystemObject) exists() bool {
	_, err := os.Stat(path.String())
	if nil != err {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (path FileSystemObject) isRoot() bool {
	return "/" == path.String()
}

func (path FileSystemObject) SearchDirectoryAbove(needle FileSystemObject) bool {
	search_directory := path.Parent().Parent()
	if path.Parent() == search_directory {
		return false
	}
	discovered, err := filepath.Glob(filepath.Join(search_directory.String(), needle.String()))
	if nil != err {
		log.Fatal(err)
	}
	return 1 <= len(discovered)
}

func (path FileSystemObject) ClimbUpwardsToFind(needle FileSystemObject) (FileSystemObject, error) {
	if path.isRoot() {
		return path, ErrRootReached
	}
	if path.SearchDirectoryAbove(needle) {
		return FileSystemObject(filepath.Join(path.Parent().String(), needle.String())), nil
	}
	return path.Parent().ClimbUpwardsToFind(needle)
}
