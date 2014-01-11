package database

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"github.com/ghthor/database/dbtesting"
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (e *RegisterAppEx) Describe(c *dbtesting.ExecutorContext) {
	var err error

	c.Impl, err = NewRegisterAppEx(c.Db)
	c.Assume(err, IsNil)

	pkgFilename := "example.pkg"
	pkgName := strings.Split(filepath.Base(pkgFilename), ".")[0]

	pkgBytes, err := ioutil.ReadFile(pkgFilename)
	c.Assume(err, IsNil)

	h := sha1.New()
	_, err = io.Copy(h, bytes.NewReader(pkgBytes))
	c.Assume(err, IsNil)

	sha1Name := hex.EncodeToString(h.Sum(nil)) + filepath.Ext(pkgFilename)

	c.SpecifyResult(App{
		Id:   AppId(1),
		Name: pkgName,
		Pkg:  File{FileId(1), sha1Name},
	})

	c.SpecifySideEffects("should insert a row in the `file` table", func() {
		conn := c.Db.MysqlDatabase().Conn
		rows, res, err := conn.Query("select * from `file`")
		c.Assume(err, IsNil)

		c.Expect(len(rows), Equals, 1)
		for _, row := range rows {
			c.Expect(row.ForceUint64(res.Map("id")), Equals, uint64(1))
			c.Expect(row.Str(res.Map("filename")), Equals, sha1Name)
		}
	})

	c.SpecifySideEffects("should save file into filesystem database", func() {
		_, err := os.Stat(filepath.Join(c.Db.Filepath(), sha1Name))
		c.Expect(os.IsNotExist(err), IsFalse)
	})

	c.SpecifySideEffects("shoud insert a row in the `app` table", func() {
		conn := c.Db.MysqlDatabase().Conn
		rows, res, err := conn.Query("select * from `app`")
		c.Assume(err, IsNil)

		c.Expect(len(rows), Equals, 1)
		for _, row := range rows {
			c.Expect(row.ForceUint64(res.Map("id")), Equals, uint64(1))
			c.Expect(row.Str(res.Map("name")), Equals, pkgName)
			c.Expect(row.ForceUint64(res.Map("pkgId")), Equals, uint64(1))
		}
	})
}

func DescribeRegisterAppExecutor(c gospec.Context) {
	pkgFile, err := testFile("example.pkg")
	c.Assume(err, IsNil)

	schemeBytes, err := ioutil.ReadFile("mysql/schema.sql")
	c.Assume(err, IsNil)

	action := RegisterApp{pkgFile}

	dbtesting.DescribeExecutor(c, action, &RegisterAppEx{}, cfg, string(schemeBytes), nil)

}
