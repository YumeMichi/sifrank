package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"sifrank/config"
)

var MysqlClient *sqlx.DB

func init() {
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
