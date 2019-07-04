package main

import "errors"

var (
	ErrUnstagedChanges  = errors.New("There are unstaged changes in your working directory")
	ErrIndexUncommitted = errors.New("There are uncommitted changes in your index")
	ErrHeadlessRepo     = errors.New("Unable to initialize in a bare repo")
)
