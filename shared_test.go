package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type GitFlowSuite struct {
	root_path FileSystemObject
	tmp_dir   FileSystemObject
}

func (suite *GitFlowSuite) SetUpSuite(c *C) {
	suite.root_path = FileSystemObject("/")
	temp_dir, err := ioutil.TempDir("", "")
	if nil != err {
		log.Fatal("Unable to build a temp directory")
	}
	suite.tmp_dir = FileSystemObject(temp_dir)
}

func (suite *GitFlowSuite) TearDownSuite(c *C) {
	os.RemoveAll(suite.tmp_dir.String())
}
