//
// Copyright 2021 YumeMichi. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// This binary provides sample code for using the gopacket TCP assembler and TCP
// stream reader.  It reads packets off the wire and reconstructs HTTP requests
// it sees, logging them.
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", config.Conf.MysqlUser, config.Conf.MysqlPassword, config.Conf.MysqlHost, config.Conf.MysqlPort, config.Conf.MysqlDb)
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
