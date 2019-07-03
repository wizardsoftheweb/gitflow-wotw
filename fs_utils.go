package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

// A simple file handle
type FileSystemObject string

var (
	// While climbing, the root path was reached
	ErrRootReached = errors.New("Root directory reached without finding file")
)

// Convert FileSystemObject to a string
func (path FileSystemObject) String() string {
	return string(path)
}

// Get the basename
func (path FileSystemObject) Base() string {
	return filepath.Base(path.String())
}

// Get the parent
func (path FileSystemObject) Parent() FileSystemObject {
	return FileSystemObject(filepath.Dir(path.String()))
}

// Check whether or not the file exists
func (path FileSystemObject) exists() bool {
	_, err := os.Stat(path.String())
	if nil != err {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Checks if a path is root
func (path FileSystemObject) isRoot() bool {
	return "/" == path.String()
}

// Searches the directory above for a specific filename
func (path FileSystemObject) SearchDirectoryAbove(needle string) bool {
	search_directory := path.Parent().Parent()
	if path.Parent() == search_directory {
		return false
	}
	discovered, err := filepath.Glob(filepath.Join(search_directory.String(), needle))
	if nil != err {
		log.Fatal(err)
	}
	return 1 <= len(discovered)
}

// Climbs up through the directory tree looking for a file
func (path FileSystemObject) ClimbUpwardsToFind(needle string) (FileSystemObject, error) {
	if path.isRoot() {
		return path, ErrRootReached
	}
	if path.SearchDirectoryAbove(needle) {
		return FileSystemObject(filepath.Join(path.Parent().Parent().String(), needle)), nil
	}
	return path.Parent().ClimbUpwardsToFind(needle)
}
