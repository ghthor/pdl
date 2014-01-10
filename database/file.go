package database

import (
	"github.com/ghthor/database"
	"github.com/ghthor/database/action"
	"github.com/ghthor/database/datatype"
	"github.com/ziutek/mymysql/mysql"
)

func init() {
	database.RegisterAction(AddFile{}, NewAddFileEx)
}

type AddFile struct {
	File datatype.FormFile
}

func (AddFile) IsValid() error { return nil }

type AddFileEx struct {
	database.DatabaseConn

	insertFileStmt mysql.Stmt
}

func (e *AddFileEx) ExecuteWith(a action.A) (interface{}, error) {
	addFile, ok := a.(AddFile)
	if !ok {
		return nil, database.ErrInvalidAction
	}

	tx, err := e.Begin()
	if err != nil {
		return nil, err
	}

	filename, err := tx.SaveFile(addFile.File)
	if err != nil {
		return nil, err
	}

	res, err := tx.Run(e.insertFileStmt, filename)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	fileId := datatype.Id(res.InsertId())

	return fileId, nil
}

func NewAddFileEx(c database.DatabaseConn) (database.Executor, error) {
	insertFileStmt, err := c.MysqlConn().Prepare("insert into `file` (filename) values (?)")
	if err != nil {
		return nil, err
	}

	return &AddFileEx{c, insertFileStmt}, nil
}
