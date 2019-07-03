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
	current_file FileSystemObject
	dummy_path   FileSystemObject
	testfile     FileSystemObject
	root_path    FileSystemObject
}

var _ = Suite(&FsUtilsSuite{})

func (suite *FsUtilsSuite) SetUpSuite(c *C) {
	suite.current_file = FileSystemObject(os.Args[0])
	suite.dummy_path = FileSystemObject(dummy_path)
	suite.testfile = FileSystemObject(filepath.Join(filepath.Dir(os.Args[0]), "testdata", "testfile"))
	suite.root_path = FileSystemObject("/")
}

func (suite *FsUtilsSuite) TestFileSystemObjectStringReturn(c *C) {
	c.Assert(suite.dummy_path.String(), Equals, dummy_path)
}

func (suite *FsUtilsSuite) TestFileSystemObjectExists(c *C) {
	c.Assert(suite.current_file.exists(), Equals, true)
	c.Assert(suite.dummy_path.exists(), Equals, false)
}

func (suite *FsUtilsSuite) TestFileSystemObjectSearchInParent(c *C) {
	c.Assert(suite.testfile.SearchDirectoryAbove("gitflow-wotw.test"), Equals, true)
	c.Assert(suite.testfile.SearchDirectoryAbove("nope"), Equals, false)
	c.Assert(suite.root_path.SearchDirectoryAbove("anything"), Equals, false)
}
