package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"
)

const (
	dummy_path = "/some/dummy/path"
)

type FsUtilsSuite struct {
	GitFlowSuite
	current_file        FileSystemObject
	dummy_path          FileSystemObject
	primary_temp_file   FileSystemObject
	secondary_temp_dir  FileSystemObject
	secondary_temp_file FileSystemObject
}

var _ = Suite(&FsUtilsSuite{})

func (suite *FsUtilsSuite) SetUpTest(c *C) {
	suite.current_file = FileSystemObject(os.Args[0])
	suite.dummy_path = FileSystemObject(dummy_path)
	primary_temp_file, err := ioutil.TempFile(suite.tmp_dir.String(), "")
	if nil != err {
		log.Fatal("Unable to build a temp directory")
	}
	secondary_temp_dir, err := ioutil.TempDir(suite.tmp_dir.String(), "")
	if nil != err {
		log.Fatal("Unable to build a temp directory")
	}
	secondary_temp_file, err := ioutil.TempFile(secondary_temp_dir, "")
	if nil != err {
		log.Fatal("Unable to build a temp directory")
	}
	suite.primary_temp_file = FileSystemObject(primary_temp_file.Name())
	suite.secondary_temp_dir = FileSystemObject(secondary_temp_dir)
	suite.secondary_temp_file = FileSystemObject(secondary_temp_file.Name())
}

func (suite *FsUtilsSuite) TestFileSystemObjectSimpleReturns(c *C) {
	c.Assert(suite.dummy_path.String(), Equals, dummy_path)
	c.Assert(suite.dummy_path.Base(), Equals, filepath.Base(dummy_path))
	c.Assert(suite.dummy_path.Parent(), Equals, FileSystemObject(filepath.Dir(dummy_path)))
	c.Assert(suite.dummy_path.isRoot(), Equals, false)
	c.Assert(suite.root_path.isRoot(), Equals, true)
}

func (suite *FsUtilsSuite) TestFileSystemObjectModeReturns(c *C) {
	c.Assert(suite.dummy_path.isFileOrDir(), Equals, FILE_MODE_ERROR)
	c.Assert(suite.dummy_path.IsDir(), Equals, false)
	c.Assert(suite.dummy_path.IsFile(), Equals, false)
	c.Assert(suite.tmp_dir.isFileOrDir(), Equals, FILE_MODE_DIRECTORY)
	c.Assert(suite.tmp_dir.IsDir(), Equals, true)
	c.Assert(suite.tmp_dir.IsFile(), Equals, false)
	c.Assert(suite.primary_temp_file.isFileOrDir(), Equals, FILE_MODE_FILE)
	c.Assert(suite.primary_temp_file.IsDir(), Equals, false)
	c.Assert(suite.primary_temp_file.IsFile(), Equals, true)
}

func (suite *FsUtilsSuite) TestFileSystemObjectDirectoryContains(c *C) {
	result, err := suite.primary_temp_file.DirectoryContains("test")
	c.Assert(err, ErrorMatches, ErrNotADirectory.Error())
	result, err = suite.tmp_dir.DirectoryContains(suite.primary_temp_file.Base())
	c.Assert(err, IsNil)
	c.Assert(result, Equals, true)
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
	discovered, _ := suite.secondary_temp_file.ClimbUpwardsToFind(suite.tmp_dir.Base())
	c.Assert(discovered, Equals, suite.tmp_dir)
	_, err := suite.root_path.ClimbUpwardsToFind("anything")
	c.Assert(err, ErrorMatches, string(ErrRootReached.Error()))
}
