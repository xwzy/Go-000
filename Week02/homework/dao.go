package main

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"log"
)

var (
	ctx context.Context
	// 此处的db未初始化
	db *sql.DB
)

func init() {
	var err error
	db, err = sql.Open("mysql",
		"root:123456@tcp(127.0.0.1:3306)/test0")
	if err != nil {
		// 数据库无法连接直接panic
		log.Panic(err)
	} else {
		log.Println("Connect MySQL success")
		ctx = context.Background()
	}
}

func GetUserName(id int) (string, error) {
	var username string
	err := db.QueryRowContext(ctx, "SELECT name FROM test0.user WHERE id=?", id).Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		return username, errors.Wrapf(err, "no user with id %d\n", id)
	case err != nil:
		return username, errors.Wrap(err, "get user id fail")
	default:
		return username, nil
	}
}
