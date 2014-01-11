package database

import (
	"errors"
	"github.com/ghthor/database"
	"github.com/ghthor/database/action"
	"github.com/ghthor/database/datatype"
	"github.com/ziutek/mymysql/mysql"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidAppName = errors.New("invalid app name")
)

type (
	AppId datatype.Id

	App struct {
		Id   AppId  `json:"id"`
		Name string `json:"name"`
		Pkg  File   `json:"pkg"`
	}
)

type RegisterApp struct {
	Pkg datatype.FormFile
}

func (a RegisterApp) IsValid() error {
	// TODO: Implement
	return nil
}

type RegisterAppEx struct {
	database.DatabaseConn

	insertFileStmt mysql.Stmt
	insertAppStmt  mysql.Stmt
}

func (e *RegisterAppEx) ExecuteWith(a action.A) (interface{}, error) {
	registerApp, ok := a.(RegisterApp)
	if !ok {
		return nil, database.ErrInvalidAction
	}

	tx, err := e.Begin()
	if err != nil {
		return nil, err
	}

	name := strings.Split(filepath.Base(registerApp.Pkg.Header.Filename), ".")[0]

	filename, err := tx.SaveFile(registerApp.Pkg)
	if err != nil {
		return nil, err
	}

	res, err := tx.Run(e.insertFileStmt, filename)
	if err != nil {
		return nil, err
	}

	fileId := FileId(res.InsertId())

	res, err = tx.Run(e.insertAppStmt, name, 1)
	if err != nil {
		return nil, err
	}

	appId := AppId(res.InsertId())

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return App{
		Id:   appId,
		Name: name,
		Pkg:  File{fileId, filename},
	}, nil
}

func NewRegisterAppEx(c database.DatabaseConn) (database.Executor, error) {
	insertFileStmt, err := c.MysqlConn().Prepare("insert into `file` (filename) values (?)")
	if err != nil {
		return nil, err
	}

	insertAppStmt, err := c.MysqlConn().Prepare("insert into `app` (name, pkgId) values (?, ?)")
	if err != nil {
		return nil, err
	}

	return &RegisterAppEx{c, insertFileStmt, insertAppStmt}, nil
}
