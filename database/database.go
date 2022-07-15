package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dextercai/db-doc-gen/model"
	_ "github.com/go-sql-driver/mysql"
)

func InitDb(config *model.DbConfig) (*DbSource, error) {
	var (
		dsn string
	)
	if *config.DbType == "mysql" {
		// https://github.com/go-sql-driver/mysql/
		// <username>:<password>@<host>:<port>/<database>
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			*config.User, *config.Password, *config.Host, *config.Port, *config.Database)
	} else {
		return nil, errors.New("不支持该数据库：" + *config.DbType)
	}
	db, err := sql.Open(*config.DbType, dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &DbSource{
		DB:     db,
		Config: *config,
	}, nil
}
