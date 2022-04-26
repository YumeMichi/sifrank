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
package main

import (
	"sifrank/config"
	"sifrank/db"
	"sifrank/sched"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func init() {
	logrus.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%time%][%lvl%]: %msg% \n",
	})

	if config.Conf.EnableMigration {
		db.InitMySQL()
		db.InitRedis()
	}
}

func main() {
	go sched.FetchNtpData()
	go sched.FetchPacketData()
	go sched.FetchRankData()

	select {}
}
