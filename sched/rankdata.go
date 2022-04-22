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
package sched

import (
	"sifrank/bot"
	"sifrank/config"
	"sifrank/db"
	"time"

	"github.com/sirupsen/logrus"
)

func FetchRankData() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	d, err := time.ParseInLocation("2006-01-02 15:04:05", config.Conf.EndTime, time.Local)
	if err != nil {
		panic(err)
	}
	endDate := d.Format("2006-01-02")

	for {
		select {
		case t := <-ticker.C:
			// 是否活动已结束
			et, _ := time.ParseInLocation("2006-01-02 15:04:05", config.Conf.EndTime, time.Local)
			if t.Unix() > et.Unix() {
				continue
			}
			// 是否活动结束当天
			h, m, s := t.Clock()
			currentDate := t.Local().Format("2006-01-02")
			var hOffset, mOffset, sOffset int
			if currentDate == endDate {
				hOffset = 13
			} else {
				hOffset = 23
			}
			mOffset = 59
			sOffset = 55
			if h == hOffset && m == mOffset && s == sOffset {
				result, err := bot.GetData()
				if err != nil || len(result) != 3 {
					logrus.Warn(err)
					return
				}
				for k, v := range result {
					prefix := time.Now().Local().Format("20060102")
					key := []byte(prefix + "_" + k)
					value := []byte(v)
					err = db.LevelDb.Put(key, value)
					if err != nil {
						logrus.Warn(err.Error())
						break
					}
				}
			}
		}
	}
}
