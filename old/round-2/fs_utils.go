package main

import (
	"errors"
	"os"
	"path/filepath"
)

// A simple file handle
type FileSystemObject string
type FileMode int

const (
	// Unknown file mode
	FILE_MODE_ERROR FileMode = iota
	// A directory
	FILE_MODE_DIRECTORY
	// A regular file
	FILE_MODE_FILE
)

var (
	// While climbing, the root path was reached
	ErrRootReached = errors.New("Root directory reached without finding file")
	// Cannot search paths that are not directories
	ErrNotADirectory = errors.New("Path must be a directory to search it")
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

// Check the file mode of the path
func (path FileSystemObject) isFileOrDir() FileMode {
	stat, err := os.Stat(path.String())
	if nil == err {
		switch mode := stat.Mode(); {
		case mode.IsDir():
			return FILE_MODE_DIRECTORY
		case mode.IsRegular():
			return FILE_MODE_FILE
		}
	}
	return FILE_MODE_ERROR
}

// Check if path is a directory
func (path FileSystemObject) IsDir() bool {
	return FILE_MODE_DIRECTORY == path.isFileOrDir()
}

// Check if path is a file
func (path FileSystemObject) IsFile() bool {
	return FILE_MODE_FILE == path.isFileOrDir()
}

// Search the path for a needle
func (path FileSystemObject) DirectoryContains(needle string) (bool, error) {
	if !path.IsDir() {
		return false, ErrNotADirectory
	}
	discovered, err := filepath.Glob(filepath.Join(path.String(), needle))
	if nil != err {
		return false, err
	}
	return 1 <= len(discovered), nil
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
	result, err := search_directory.DirectoryContains(needle)
	if nil != err {
		return false
	}
	return result
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
