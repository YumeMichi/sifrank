//
// Copyright (C) 2021 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package db

import (
	"fmt"
	"sifrank/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var MysqlClient *sqlx.DB

func InitMySQL() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", config.Conf.Mysql.MysqlUser, config.Conf.Mysql.MysqlPassword,
		config.Conf.Mysql.MysqlHost, config.Conf.Mysql.MysqlPort, config.Conf.Mysql.MysqlDb)
	client, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(5)
	err = client.Ping()
	if err != nil {
		panic(err.Error())
	}
	MysqlClient = client
}
