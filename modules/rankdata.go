//
// Copyright (C) 2022 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package modules

import (
	"sifrank/bot"
	"sifrank/config"
	"sifrank/db"
	"sifrank/xclog"
	"time"
)

func RankDataTicker() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	d, err := time.ParseInLocation("2006-01-02 15:04:05", config.Conf.Event.EndTime, time.Local)
	if err != nil {
		panic(err)
	}
	endDate := d.Format("2006-01-02")

	for {
		select {
		case t := <-ticker.C:
			// 是否活动已结束
			et, _ := time.ParseInLocation("2006-01-02 15:04:05", config.Conf.Event.EndTime, time.Local)
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
					xclog.Warn(err)
					return
				}
				for k, v := range result {
					prefix := time.Now().Local().Format("20060102")
					key := []byte(prefix + "_" + k)
					value := []byte(v)
					err = db.LevelDb.Put(key, value)
					if err != nil {
						xclog.Warn(err.Error())
						break
					}
				}
			}
		}
	}
}
