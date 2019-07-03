package main

import (
	"os"

	. "gopkg.in/check.v1"
)

const (
	dummy_path = "/some/dummy/path"
)

type FsUtilsSuite struct {
	current_file FileSystemObject
	dummy_path   FileSystemObject
}

var _ = Suite(&FsUtilsSuite{})

func (suite *FsUtilsSuite) SetUpSuite(c *C) {
	suite.current_file = FileSystemObject(os.Args[0])
	suite.dummy_path = FileSystemObject(dummy_path)
}

func (suite *FsUtilsSuite) TestFileSystemObjectStringReturn(c *C) {
	c.Assert(suite.dummy_path.String(), Equals, dummy_path)
}

func (suite *FsUtilsSuite) TestFileSystemObjectExists(c *C) {
	c.Assert(suite.current_file.exists(), Equals, true)
}
