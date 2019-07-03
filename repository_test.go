package main

import (
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"
)

type RepositorySuite struct {
	GitFlowSuite
	repo        *Repository
	dot_git_dir FileSystemObject
}

var _ = Suite(&RepositorySuite{})

func (suite *RepositorySuite) SetUpTest(c *C) {
	suite.repo = &Repository{}
	suite.dot_git_dir = FileSystemObject(filepath.Join(suite.tmp_dir.String(), ".git"))
	os.Mkdir(suite.dot_git_dir.String(), 0777)
}

func (suite *RepositorySuite) TestDiscoverDotDir(c *C) {
	result, err := suite.repo.discoverDotDir(suite.root_path)
	c.Assert(err, ErrorMatches, ErrNotARepo.Error())
	c.Assert(result, Equals, suite.root_path)
	result, err = suite.repo.discoverDotDir(suite.tmp_dir)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, suite.dot_git_dir)
}
