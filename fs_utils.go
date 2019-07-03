package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type FileSystemObject string

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

func (path FileSystemObject) SearchInParent(needle FileSystemObject) bool {
	discovered, err := filepath.Glob(fmt.Sprintf("%s/%s", path.Parent().String(), needle.String()))
	if nil != err {
		log.Fatal(err)
	}
	return 1 <= len(discovered)
}
