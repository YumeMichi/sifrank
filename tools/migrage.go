//
// Copyright (C) 2022 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package tools

import (
	"fmt"
	"sifrank/db"
	"sifrank/xclog"
	"strconv"
	"strings"
)

type DayRankData struct {
	Id       int    `db:"id"`
	Rank     string `db:"rank"`
	RankCode string `db:"rank_code"`
	Score    int    `db:"score"`
	DataDate string `db:"data_date"`
	DataTime string `db:"data_time"`
}

func MigrateFromMySQLToLevelDB() {
	var data []DayRankData
	err := db.MysqlClient.Select(&data, "SELECT * FROM day_rank_data")
	if err != nil {
		xclog.Warn(err.Error())
		return
	}

	for _, v := range data {
		prefix := strings.ReplaceAll(v.DataDate, "-", "")
		key := []byte(prefix + "_" + v.Rank)
		value := []byte(strconv.Itoa(v.Score))
		err = db.LevelDb.Put([]byte(key), []byte(value))
		err = db.LevelDb.Put(key, value)
		if err != nil {
			xclog.Warn(err.Error())
		}
		fmt.Println("==== Put ", string(key), " ====")
	}

	list := db.LevelDb.List()
	fmt.Println(list)
}
