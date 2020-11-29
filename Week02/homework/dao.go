package main

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"log"
)

type DAO struct {
	ctx context.Context
	db  *sql.DB
}

func DAOFactory() *DAO {
	dao := new(DAO)
	var err error
	dao.db, err = sql.Open("mysql",
		"root:123456@tcp(127.0.0.1:3306)/test0")
	if err != nil {
		// 数据库无法连接直接panic
		log.Panic(err)
		return nil
	} else {
		log.Println("Connect MySQL success")
		dao.ctx = context.Background()
	}
	return dao
}

func (dao *DAO) GetUserName(id int) (string, error) {
	var username string
	err := dao.db.QueryRowContext(dao.ctx, "SELECT name FROM test0.user WHERE id=?", id).Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		return username, errors.WithMessagef(err, "no user with id %d\n", id)
	case err != nil:
		return username, errors.Wrap(err, "get user id fail")
	default:
		return username, nil
	}
}
