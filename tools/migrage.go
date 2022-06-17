//
// Copyright (C) 2022 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package tools

import (
	"fmt"
	"sifrank/db"
	"sifrank/model"
	"sifrank/xclog"
	"strconv"
	"strings"
)

func MigrateFromMySQLToLevelDB() {
	var data []model.DayRankData
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
