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
	"mime/multipart"
	"os"
	"path/filepath"
)

func (e *AddFileEx) Describe(c *dbtesting.ExecutorContext) {
	var err error

	c.Impl, err = NewAddFileEx(c.Db)
	c.Assume(err, IsNil)

	expectedBytes, err := ioutil.ReadFile("image_test.png")
	c.Assume(err, IsNil)

	h := sha1.New()
	_, err = io.Copy(h, bytes.NewReader(expectedBytes))
	c.Assume(err, IsNil)

	sha1Name := hex.EncodeToString(h.Sum(nil)) + filepath.Ext("image_test.png")

	c.SpecifyResult(File{FileId(1), sha1Name})

	c.SpecifySideEffects("should insert a row into the `file` table", func() {
		conn := c.Db.MysqlDatabase().Conn
		rows, res, err := conn.Query("select * from `file`")
		c.Assume(err, IsNil)
		c.Assume(len(rows), Equals, 1)
		for _, row := range rows {
			c.Expect(row.ForceUint64(res.Map("id")), Equals, uint64(1))
			c.Expect(row.Str(res.Map("filename")), Equals, sha1Name)
		}
	})

	c.SpecifySideEffects("should save file into filesystem database", func() {
		var filenames []string
		walker := func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				filenames = append(filenames, path)
				bytes, err := ioutil.ReadFile(path)
				c.Assume(err, IsNil)

				c.Expect(string(bytes), Equals, string(expectedBytes))
				c.Expect(filepath.Base(path), Equals, sha1Name)
			}
			return nil
		}

		c.Assume(filepath.Walk(c.Db.Filepath(), walker), IsNil)
		c.Assume(len(filenames), Equals, 1)
	})

}

func testFile(filepathStr string) (datatype.FormFile, error) {
	file, err := os.Open(filepathStr)
	if err != nil {
		return datatype.FormFile{}, err
	}

	return datatype.FormFile{
		File: file,
		Header: &multipart.FileHeader{
			Filename: filepath.Base(filepathStr),
		},
	}, nil
}

func DescribeAddFileExecutor(c gospec.Context) {
	testImage, err := testFile("image_test.png")
	c.Assume(err, IsNil)

	schemeBytes, err := ioutil.ReadFile("mysql/schema.sql")
	c.Assume(err, IsNil)

	addFile := AddFile{testImage}

	dbtesting.DescribeExecutor(c, addFile, &AddFileEx{}, cfg, string(schemeBytes), nil)
}
