package main

import (
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"
)

const (
	dummy_path = "/some/dummy/path"
)

type FsUtilsSuite struct {
	current_file     FileSystemObject
	dummy_path       FileSystemObject
	test_file        FileSystemObject
	root_path        FileSystemObject
	deeper_test_file FileSystemObject
}

var _ = Suite(&FsUtilsSuite{})

func (suite *FsUtilsSuite) SetUpSuite(c *C) {
	suite.current_file = FileSystemObject(os.Args[0])
	suite.dummy_path = FileSystemObject(dummy_path)
	suite.test_file = FileSystemObject(filepath.Join(filepath.Dir(os.Args[0]), "testdata", "test_file"))
	suite.root_path = FileSystemObject("/")
	suite.deeper_test_file = FileSystemObject(filepath.Join(filepath.Dir(os.Args[0]), "testdata", "deeper", "file"))
}

func (suite *FsUtilsSuite) TestFileSystemObjectStringReturn(c *C) {
	c.Assert(suite.dummy_path.String(), Equals, dummy_path)
}

func (suite *FsUtilsSuite) TestFileSystemObjectExists(c *C) {
	c.Assert(suite.current_file.exists(), Equals, true)
	c.Assert(suite.dummy_path.exists(), Equals, false)
}

func (suite *FsUtilsSuite) TestFileSystemObjectSearchInParent(c *C) {
	c.Assert(suite.deeper_test_file.SearchDirectoryAbove("test_file"), Equals, true)
	c.Assert(suite.deeper_test_file.SearchDirectoryAbove("nope"), Equals, false)
	c.Assert(suite.root_path.SearchDirectoryAbove("anything"), Equals, false)
}

func (suite *FsUtilsSuite) TestFileSystemObjectClimbUpwardsToFind(c *C) {
	discovered, _ := suite.deeper_test_file.ClimbUpwardsToFind("test_file")
	c.Assert(discovered, Equals, suite.test_file)
	_, err := suite.root_path.ClimbUpwardsToFind("anything")
	c.Assert(err, ErrorMatches, string(ErrRootReached.Error()))
}
