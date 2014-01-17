package database

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"github.com/ghthor/database/datatype"
	"github.com/ghthor/database/dbtesting"
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func pkgFiles() (map[string]datatype.FormFile, error) {
	pkgs := map[string]datatype.FormFile{
		"example.pkg":   datatype.FormFile{},
		"example.pkg.1": datatype.FormFile{},
		"example.pkg.2": datatype.FormFile{},
		"another.pkg":   datatype.FormFile{},
	}

	for filename, _ := range pkgs {
		formFile, err := testFile(filename)
		if err != nil {
			return nil, err
		}

		pkgs[filename] = formFile
	}
	pkgs["example.pkg.1"].Header.Filename = "example.pkg"
	pkgs["example.pkg.2"].Header.Filename = "example.pkg"

	return pkgs, nil
}

func (e *InstallAppEx) Describe(c *dbtesting.ExecutorContext) {
	var err error

	c.Impl, err = NewInstallAppEx(c.Db)
	c.Assume(err, IsNil)

	c.Specify("after executing without error", func() {
		pkgFilename := "example.pkg"
		pkgName := strings.Split(filepath.Base(pkgFilename), ".")[0]

		pkgBytes, err := ioutil.ReadFile(pkgFilename)
		c.Assume(err, IsNil)

		h := sha1.New()
		_, err = io.Copy(h, bytes.NewReader(pkgBytes))
		c.Assume(err, IsNil)

		sha1Name := hex.EncodeToString(h.Sum(nil)) + filepath.Ext(pkgFilename)

		c.SpecifySideEffects("will have inserted a row in the `file` table", func() {
			conn := c.Db.MysqlDatabase().Conn
			rows, res, err := conn.Query("select * from `file`")
			c.Assume(err, IsNil)

			c.Expect(len(rows), Equals, 1)
			for _, row := range rows {
				c.Expect(row.ForceUint64(res.Map("id")), Equals, uint64(1))
				c.Expect(row.Str(res.Map("filename")), Equals, sha1Name)
			}
		})

		c.SpecifySideEffects("will have saved pkg into filesystem database", func() {
			_, err := os.Stat(filepath.Join(c.Db.Filepath(), sha1Name))
			c.Expect(os.IsNotExist(err), IsFalse)
		})

		c.SpecifySideEffects("will have inserted a row in the `app` table", func() {
			conn := c.Db.MysqlDatabase().Conn
			rows, res, err := conn.Query("select id, name, pkgId from `app`")
			c.Assume(err, IsNil)

			c.Expect(len(rows), Equals, 1)
			for _, row := range rows {
				c.Expect(row.ForceUint64(res.Map("id")), Equals, uint64(1))
				c.Expect(row.Str(res.Map("name")), Equals, pkgName)
				c.Expect(row.ForceUint64(res.Map("pkgId")), Equals, uint64(1))
			}
		})

		c.SpecifyResult(App{
			Id:   AppId(1),
			Name: pkgName,
			Pkg:  File{FileId(1), sha1Name},
		})
	})

	c.Specify("will fail because an app with the name already exists and", func() {
		// Install the App
		ex, err := NewInstallAppEx(c.Db)
		c.Assume(err, IsNil)

		_, err = ex.ExecuteWith(c.Input)
		c.Assume(err, IsNil)

		// Setup to Install the app again with a different pkg file
		// This is to verify that the installed app's pkg isn't overwritten
		pkgFilename := "example.pkg.1"
		pkgFile, err := testFile(pkgFilename)
		c.Assume(err, IsNil)
		// Make sure the Filename's are the same as far as the Executor is concerned
		pkgFile.Header.Filename = "example.pkg"
		c.Input = InstallApp{pkgFile}

		pkgBytes, err := ioutil.ReadFile(pkgFilename)
		c.Assume(err, IsNil)

		h := sha1.New()
		_, err = io.Copy(h, bytes.NewReader(pkgBytes))
		c.Assume(err, IsNil)

		sha1Name := hex.EncodeToString(h.Sum(nil)) + filepath.Ext(pkgFilename)

		c.SpecifySideEffects("will not insert a row in the `file` table", func() {
			conn := c.Db.MysqlDatabase().Conn
			rows, _, err := conn.Query("select * from `file`")
			c.Assume(err, IsNil)

			c.Expect(len(rows), Equals, 1)
		})

		c.SpecifySideEffects("will not save pkg into filesystem database", func() {
			var filenames []string
			walker := func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					filenames = append(filenames, path)
					bytes, err := ioutil.ReadFile(path)
					c.Assume(err, IsNil)

					// Check to make sure this file IS NOT example.pkg.1
					c.Expect(string(bytes), Not(Equals), string(pkgBytes))
					c.Expect(filepath.Base(path), Not(Equals), sha1Name)
				}
				return nil
			}

			c.Assume(filepath.Walk(c.Db.Filepath(), walker), IsNil)
			c.Expect(len(filenames), Equals, 1)
		})

		c.SpecifySideEffects("will not insert a row in the `app` table", func() {
			conn := c.Db.MysqlDatabase().Conn
			rows, _, err := conn.Query("select * from `app`")
			c.Assume(err, IsNil)

			c.Expect(len(rows), Equals, 1)
		})

		c.Expect(c.Res, IsNil)
		c.Expect(c.Err, Equals, ErrAppAlreadyExists)
	})
}

func (e *UpdateAppEx) Describe(c *dbtesting.ExecutorContext) {
	var err error

	c.Impl, err = NewUpdateAppEx(c.Db)
	c.Assume(err, IsNil)

	pkgs, err := pkgFiles()
	c.Assume(err, IsNil)

	// Install example.pkg
	installAppEx, err := NewInstallAppEx(c.Db)
	c.Assume(err, IsNil)

	_, err = installAppEx.ExecuteWith(InstallApp{pkgs["example.pkg"]})
	c.Assume(err, IsNil)

	c.Specify("after executing without an error", func() {
		pkgBytes, err := ioutil.ReadFile("example.pkg.1")
		c.Assume(err, IsNil)

		h := sha1.New()
		_, err = io.Copy(h, bytes.NewReader(pkgBytes))
		c.Assume(err, IsNil)

		sha1Name := hex.EncodeToString(h.Sum(nil)) + filepath.Ext("example.pkg")

		c.SpecifySideEffects("will save the pkg to the filesystem", func() {
			_, err := os.Stat(filepath.Join(c.Db.Filepath(), sha1Name))
			c.Expect(os.IsNotExist(err), IsFalse)
		})

		// Because the pkg file is different
		c.SpecifySideEffects("will insert a new row in the `file` table for the pkg", func() {
			conn := c.Db.MysqlDatabase().Conn
			rows, res, err := conn.Query("select * from `file`")
			c.Assume(err, IsNil)

			c.Expect(len(rows), Equals, 2)
			for _, row := range rows[1:] {
				c.Expect(row.ForceUint64(res.Map("id")), Equals, uint64(2))
				c.Expect(row.Str(res.Map("filename")), Equals, sha1Name)
			}
		})

		// App Updates will be stored as another app with the same name
		// When loading an app by name, the app with the most recent install date will be used
		c.SpecifySideEffects("will insert a new row in the `app` table", func() {
			conn := c.Db.MysqlDatabase().Conn
			rows, res, err := conn.Query("select id, name, pkgId from `app`")
			c.Assume(err, IsNil)

			c.Expect(len(rows), Equals, 2)
			for _, row := range rows[1:] {
				c.Expect(row.ForceUint64(res.Map("id")), Equals, uint64(2))
				c.Expect(row.Str(res.Map("name")), Equals, "example")
				c.Expect(row.ForceUint64(res.Map("pkgId")), Equals, uint64(2))
			}
		})

		c.SpecifyResult(App{
			Id:   AppId(2),
			Name: "example",
			Pkg:  File{FileId(2), sha1Name},
		})
	})

	notSaved := func() {}
	noFileInserted := func() {}
	noAppInserted := func() {}

	c.Specify("will fail because the app has not been installed and", func() {
		c.SpecifySideEffects("will not save the pkg to the filesystem", notSaved)

		c.SpecifySideEffects("will not insert a new row in the `file` table", noFileInserted)

		c.SpecifySideEffects("will not insert a new row in the `app` table", noAppInserted)
	})

	c.Specify("will fail because the pkg file is not different and", func() {
		c.SpecifySideEffects("will not save the pkg to the filesystem", notSaved)

		c.SpecifySideEffects("will not insert a new row in the `file` table", noFileInserted)

		c.SpecifySideEffects("will not insert a new row in the `app` table", noAppInserted)
	})
}

func DescribeExecutors(c gospec.Context) {
	pkgs, err := pkgFiles()
	c.Assume(err, IsNil)

	schemeBytes, err := ioutil.ReadFile("mysql/schema.sql")
	c.Assume(err, IsNil)

	dbtesting.DescribeExecutor(c, InstallApp{pkgs["example.pkg"]}, &InstallAppEx{}, cfg, string(schemeBytes), nil)
	dbtesting.DescribeExecutor(c, UpdateApp{pkgs["example.pkg.1"]}, &UpdateAppEx{}, cfg, string(schemeBytes), nil)
}
