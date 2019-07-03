package main

import (
	"os"
)

type FileSystemObject string

func (path FileSystemObject) String() string {
	return string(path)
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
