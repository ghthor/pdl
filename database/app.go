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

var (
	ErrAppAlreadyExists = errors.New("app already exists")
)

type (
	AppId datatype.Id

	App struct {
		Id   AppId  `json:"id"`
		Name string `json:"name"`
		Pkg  File   `json:"pkg"`
	}
)

type InstallApp struct {
	Pkg datatype.FormFile
}

func (a InstallApp) IsValid() error {
	// TODO: Implement
	return nil
}

type InstallAppEx struct {
	database.DatabaseConn

	selectAppStmt  mysql.Stmt
	insertFileStmt mysql.Stmt
	insertAppStmt  mysql.Stmt
}

func (e *InstallAppEx) ExecuteWith(a action.A) (interface{}, error) {
	installApp, ok := a.(InstallApp)
	if !ok {
		return nil, database.ErrInvalidAction
	}

	name := strings.Split(filepath.Base(installApp.Pkg.Header.Filename), ".")[0]
	rows, _, err := e.selectAppStmt.Exec(name)
	if err != nil {
		return nil, err
	}

	if len(rows) != 0 {
		return nil, ErrAppAlreadyExists
	}

	tx, err := e.Begin()
	if err != nil {
		return nil, err
	}

	filename, err := tx.SaveFile(installApp.Pkg)
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

func NewInstallAppEx(c database.DatabaseConn) (database.Executor, error) {
	selectAppStmt, err := c.MysqlConn().Prepare("select (name) from `app` where name = ?")
	if err != nil {
		return nil, err
	}

	insertFileStmt, err := c.MysqlConn().Prepare("insert into `file` (filename) values (?)")
	if err != nil {
		return nil, err
	}

	insertAppStmt, err := c.MysqlConn().Prepare("insert into `app` (name, pkgId) values (?, ?)")
	if err != nil {
		return nil, err
	}

	return &InstallAppEx{c, selectAppStmt, insertFileStmt, insertAppStmt}, nil
}
