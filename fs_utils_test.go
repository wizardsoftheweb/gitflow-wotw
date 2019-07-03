package main

import (
	"io/ioutil"
	"log"
	"os"

	. "gopkg.in/check.v1"
)

const (
	dummy_path = "/some/dummy/path"
)

type FsUtilsSuite struct {
	root_path           FileSystemObject
	current_file        FileSystemObject
	dummy_path          FileSystemObject
	primary_temp_dir    FileSystemObject
	primary_temp_file   FileSystemObject
	secondary_temp_dir  FileSystemObject
	secondary_temp_file FileSystemObject
}

var _ = Suite(&FsUtilsSuite{})

func (suite *FsUtilsSuite) SetUpSuite(c *C) {
	suite.root_path = FileSystemObject("/")
	suite.current_file = FileSystemObject(os.Args[0])
	suite.dummy_path = FileSystemObject(dummy_path)
	primary_temp_dir, err := ioutil.TempDir("", "")
	if nil != err {
		log.Fatal("Unable to build a temp directory")
	}
	primary_temp_file, err := ioutil.TempFile(primary_temp_dir, "")
	if nil != err {
		log.Fatal("Unable to build a temp directory")
	}
	secondary_temp_dir, err := ioutil.TempDir(primary_temp_dir, "")
	if nil != err {
		log.Fatal("Unable to build a temp directory")
	}
	secondary_temp_file, err := ioutil.TempFile(secondary_temp_dir, "")
	if nil != err {
		log.Fatal("Unable to build a temp directory")
	}
	suite.primary_temp_dir = FileSystemObject(primary_temp_dir)
	suite.primary_temp_file = FileSystemObject(primary_temp_file.Name())
	suite.secondary_temp_dir = FileSystemObject(secondary_temp_dir)
	suite.secondary_temp_file = FileSystemObject(secondary_temp_file.Name())
}

func (suite *FsUtilsSuite) TearDownSuite(c *C) {
	os.RemoveAll(suite.primary_temp_dir.String())
}

func (suite *FsUtilsSuite) TestFileSystemObjectStringReturn(c *C) {
	c.Assert(suite.dummy_path.String(), Equals, dummy_path)
}

func (suite *FsUtilsSuite) TestFileSystemObjectExists(c *C) {
	c.Assert(suite.current_file.exists(), Equals, true)
	c.Assert(suite.dummy_path.exists(), Equals, false)
}

func (suite *FsUtilsSuite) TestFileSystemObjectSearchInParent(c *C) {
	c.Assert(suite.secondary_temp_file.SearchDirectoryAbove(suite.primary_temp_file.Base()), Equals, true)
	c.Assert(suite.secondary_temp_file.SearchDirectoryAbove("nope"), Equals, false)
	c.Assert(suite.root_path.SearchDirectoryAbove("anything"), Equals, false)
}

func (suite *FsUtilsSuite) TestFileSystemObjectClimbUpwardsToFind(c *C) {
	discovered, _ := suite.secondary_temp_file.ClimbUpwardsToFind(suite.primary_temp_dir.Base())
	c.Assert(discovered, Equals, suite.primary_temp_dir)
	_, err := suite.root_path.ClimbUpwardsToFind("anything")
	c.Assert(err, ErrorMatches, string(ErrRootReached.Error()))
}
