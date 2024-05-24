// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package readdirplusfs_test

import (
	"os"
	"path"
	"testing"

	"github.com/jacobsa/fuse/fusetesting"
	"github.com/jacobsa/fuse/samples"
	"github.com/jacobsa/fuse/samples/readdirplusfs"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
)

func TestReadDirPlusFS(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////

type ReadDirPlusFSTest struct {
	samples.SampleTest
}

func init() { RegisterTestSuite(&ReadDirPlusFSTest{}) }

func (t *ReadDirPlusFSTest) SetUp(ti *TestInfo) {
	var err error

	t.Server, err = readdirplusfs.NewReadDirPlusFS(&t.Clock)
	AssertEq(nil, err)

	t.SampleTest.SetUp(ti)
}

////////////////////////////////////////////////////////////////////////
// Test functions
////////////////////////////////////////////////////////////////////////

func (t *ReadDirPlusFSTest) ReadDir_Root() {
	entries, err := fusetesting.ReadDirPicky(t.Dir)

	AssertEq(nil, err)
	AssertEq(2, len(entries))
	var fi os.FileInfo

	// dir
	fi = entries[0]
	ExpectEq("dir", fi.Name())
	ExpectEq(0, fi.Size())
	ExpectEq(os.ModeDir|0555, fi.Mode())
	ExpectEq(0, t.Clock.Now().Sub(fi.ModTime()), "ModTime: %v", fi.ModTime())
	ExpectTrue(fi.IsDir())

	// hello
	fi = entries[1]
	ExpectEq("hello", fi.Name())
	ExpectEq(len("Hello, world!"), fi.Size())
	ExpectEq(0444, fi.Mode())
	ExpectEq(0, t.Clock.Now().Sub(fi.ModTime()), "ModTime: %v", fi.ModTime())
	ExpectFalse(fi.IsDir())
}

func (t *ReadDirPlusFSTest) ReadDir_Dir() {
	entries, err := fusetesting.ReadDirPicky(path.Join(t.Dir, "dir"))

	AssertEq(nil, err)
	AssertEq(1, len(entries))
	var fi os.FileInfo

	// world
	fi = entries[0]
	ExpectEq("world", fi.Name())
	ExpectEq(len("Hello, world!"), fi.Size())
	ExpectEq(0444, fi.Mode())
	ExpectEq(0, t.Clock.Now().Sub(fi.ModTime()), "ModTime: %v", fi.ModTime())
	ExpectFalse(fi.IsDir())
}

func (t *ReadDirPlusFSTest) ReadDir_NonExistent() {
	_, err := fusetesting.ReadDirPicky(path.Join(t.Dir, "foobar"))

	AssertNe(nil, err)
	ExpectThat(err, Error(HasSubstr("no such file")))
}
