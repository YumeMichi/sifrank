//
// Copyright (C) 2021-2022 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package main

import (
	"sifrank/config"
	"sifrank/db"
	"sifrank/sched"
	"sifrank/xclog"
)

func init() {
	if config.Conf.EnableMigration {
		db.InitMySQL()
		db.InitRedis()
	}

	xclog.Init(config.Conf.Log.LogDir, "", config.Conf.Log.LogLevel, config.Conf.Log.LogSave)
}

func main() {
	go sched.FetchNtpData()
	go sched.FetchPacketData()
	go sched.FetchRankData()

	select {}
}
