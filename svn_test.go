package svn

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBasic(t *testing.T) {
	r, err := Open("testdata/sample")

	if err != nil {
		t.Fatal(err)
	}

	defer r.Close()

	rev, err := r.LatestRevision()

	if err != nil {
		t.Fatal(err)
	}

	if rev <= 0 {
		t.Errorf("Latest revision should be >= 0, but it is %d", rev)
	}

	c, err := r.CommitInfo(1)

	if err != nil {
		t.Fatal(err)
	}

	if c.Author != "lz" {
		t.Fatalf("Wrong author: '%s'", c.Author)
	}

	commits, err := r.Commits(1, 2)

	if err != nil {
		t.Fatal(err)
	}

	commits_count := len(commits)

	if commits_count != 2 {
		t.Error("it should return 2 commits, got", commits_count)
	}

	for _, commit := range commits {
		if commit.Author != "lz" {
			t.Errorf("Wrong author: '%s'", commit.Author)
		}
	}

	rev, err = r.LastPathRev("trunk/Makefile", 6)

	if err != nil {
		t.Error("failed to get rev", err)
	} else if rev != 5 {
		t.Error("Extected rev 5, got", rev)
	}

	commits, err = r.History("trunk", 0, rev, 2)

	if err != nil {
		t.Fatal(err)
	}

	if len(commits) != 2 {
		t.Error("#History should return 2 commits")
	}

	entries, err := r.Tree("trunk", 5)

	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 2 {
		t.Error("Only two entries should be in trunk/ folder at rev 5")
	}

	if size, err := r.FileSize("trunk/Makefile", 6); err != nil {
		t.Fatal(err)
	} else {
		exp := int64(1279)

		if size != exp {
			t.Error("Wrong file size, expected", exp, "got", size)
		}
	}

	mimeExp := "application/octet-stream"

	if mime, err := r.MimeType("trunk/images/play.png", 7); err != nil {
		t.Fatal(err)
	} else {
		if mimeExp != mime {
			t.Error("Wrong file mime type, expected", mimeExp, "got", mime)
		}
	}

	if reader, err := r.FileContent("trunk/TODO", 6); err != nil {
		t.Fatal(err)
	} else {
		data, e := ioutil.ReadAll(reader)

		if e != nil {
			t.Fatal(e)
		}

		if string(data) != "Readme\n" {
			t.Error("Wrong trunk/TODO content", string(data))
		}
	}

	ci, err := r.Changeset(5, false)

	if err != nil {
		t.Fatal(err)
	}

	diff := `Index: trunk/Makefile
===================================================================
--- trunk/Makefile	(revision 4)
+++ trunk/Makefile	(revision 5)
@@ -1,3 +1,4 @@
+# Make file to build newbc project
 export GOPATH := $(CURDIR)
 export LIBGIT_INSTALL_PREFIX := $(CURDIR)/vendor/libgit2_bin
 export LIBGIT_SRC_PATH := $(CURDIR)/vendor/libgit2
`
	if ci.ChangedPaths["trunk/Makefile"].Diff != diff {
		t.Errorf("Wrong file diff:\n%s\nExpected:\n%s", ci.ChangedPaths["trunk/Makefile"].Diff, diff)
	}

	ci, err = r.Changeset(6, false)

	if err != nil {
		t.Fatal(err)
	}

	diff = `Index: trunk/TODO
===================================================================
--- trunk/TODO	(revision 0)
+++ trunk/TODO	(revision 6)
@@ -0,0 +1 @@
+Readme
`
	if ci.ChangedPaths["trunk/TODO"].Diff != diff {
		t.Errorf("Wrong file diff:\n%s\nExpected:\n%s", ci.ChangedPaths["trunk/TODO"].Diff, diff)
	}

	ci, err = r.Changeset(9, true)

	if err != nil {
		t.Fatal(err)
	}

	if len(ci.ChangedPaths["trunk/Makefile"].Diff) > 0 {
		t.Errorf("Empty diff expected, but got \n%s", ci.ChangedPaths["trunk/Makefile"].Diff)
	}

	log := "White space change"

	if ci.Commit.Log != log {
		t.Errorf("Bad commit log message, expected '%s', got '%s'", log, ci.Commit.Log)
	}

	value, err := r.PropGet("trunk/img", 10, "svn:special")

	if err != nil {
		t.Error("Can not get prop", err.Error())
	}

	if value != "*" {
		t.Error("Bad prop value for svn:special, got:", value)
	}

	props, err := r.PropList("trunk/img", 10)

	if err != nil {
		t.Error(err)
	}

	if len(props) != 1 {
		t.Error("number of properties should be 1")
	}

	testPropKey := "svn:special"
	testPropVal := "*"
	if props[testPropKey] != testPropVal {
		t.Errorf("expected to get:%q, got: %q", testPropVal, props[testPropKey])
	}
}

func TestCreateRepo(t *testing.T) {
	path := filepath.Join(os.TempDir(), "svn-repo")
	os.RemoveAll(path)
	err := Create(path)

	if err != nil {
		t.Fatal(err)
	}
}
