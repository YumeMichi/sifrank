//
// Copyright 2022 YumeMichi. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// This binary provides sample code for using the gopacket TCP assembler and TCP
// stream reader.  It reads packets off the wire and reconstructs HTTP requests
// it sees, logging them.
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
