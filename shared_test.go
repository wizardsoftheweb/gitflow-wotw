package main

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type GitFlowSuite struct {
	root_path FileSystemObject
}

func (suite *GitFlowSuite) SetUpSuite(c *C) {
	suite.root_path = FileSystemObject("/")
}
