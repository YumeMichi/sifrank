//
// Copyright (C) 2021-2022 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package main

import (
	"sifrank/config"
	"sifrank/modules"
	"sifrank/xclog"
)

func init() {
	xclog.Init(config.Conf.Log.LogDir, "", config.Conf.Log.LogLevel, config.Conf.Log.LogSave)
}

func main() {
	go modules.SyncNtpDate()
	go modules.CapPackets()
	go modules.RankDataTicker()

	select {}
}
