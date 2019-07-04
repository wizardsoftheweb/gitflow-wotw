package main

import "errors"

var (
	ErrAlreadyInitialized = errors.New("Gitflow is already initialized; use -f to force reinit")
	ErrUnstagedChanges    = errors.New("There are unstaged changes in your working directory")
	ErrIndexUncommitted   = errors.New("There are uncommitted changes in your index")
	ErrHeadlessRepo       = errors.New("Unable to initialize in a bare repo")
	ErrProdDoesntExist    = errors.New("The production branch does not exist")
)
