//
// Copyright (C) 2022 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package cmd

import (
	"math"
	"sifrank/config"
	"sifrank/db"
	"sifrank/utils"
	"sifrank/xclog"
	"strconv"
	"time"

	"github.com/fogleman/gg"
)

var (
	x_offset = config.Conf.Dq.DqXOffset
	x_step   = config.Conf.Dq.DqXStep
	x_extra  = config.Conf.Dq.DqXExtra
	y_offset = config.Conf.Dq.DqYOffset
	y_step   = config.Conf.Dq.DqYStep
)

func GenDayRankPic() (string, error) {
	startDate, err := time.ParseInLocation("2006-01-02 15:04:05", config.Conf.Event.StartTime, time.Local)
	if err != nil {
		return "", err
	}
	endDate, err := time.ParseInLocation("2006-01-02 15:04:05", config.Conf.Event.EndTime, time.Local)
	if err != nil {
		return "", err
	}
	now := time.Now().Format("2006-01-02")
	isEnd := now == endDate.Format("2006-01-02")
	fileName := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if isEnd {
		fileName = now
	}
	savePath := config.Conf.Dq.DqOutputDir + "/" + fileName + ".png"
	if !isEnd && utils.PathExists(config.Conf.Dq.DqOutputDir+"/"+fileName+".png") {
		return savePath, nil
	}
	img, err := gg.LoadImage(config.Conf.Dq.DqBaseImage)
	if err != nil {
		return "", err
	}
	dc := gg.NewContextForImage(img)
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace(config.Conf.Dq.DqFontName, config.Conf.Dq.DqFontSize); err != nil {
		return "", err
	}
	// 档线大标题
	dc.DrawString(config.Conf.Dq.DqTitle, config.Conf.Dq.DqTitleXOffset, config.Conf.Dq.DqTitleYOffset)
	if config.Conf.Dq.DqSubtitle != "" {
		dc.SetRGB(1, 0, 0)
		dc.DrawString(config.Conf.Dq.DqSubtitle, config.Conf.Dq.DqSubtitleXOffset, config.Conf.Dq.DqSubtitleYOffset)
		dc.SetRGB(0, 0, 0)
	}
	// 档线小标题
	dc.DrawString("120", float64(x_offset), float64(y_offset))
	dc.DrawString("700", float64(x_offset), float64(y_offset+y_step*1))
	dc.DrawString("2300", float64(x_offset), float64(y_offset+y_step*2))
	// 档线数据
	dayDiff := int(math.Ceil(endDate.Sub(startDate).Hours()/24)) + 1
	if err != nil {
		return "", err
	}
	for i := 1; i <= dayDiff; i++ {
		timeDiff, _ := time.ParseDuration(strconv.Itoa(24*(i-1)) + "h")
		rankDate := startDate.Add(timeDiff).Format("20060102")
		list := db.LevelDb.ListPrefix([]byte(rankDate))
		xclog.Debug(rankDate)
		if len(list) == 0 {
			break
		}
		r1 := list[rankDate+"_ranking_1"]
		r2 := list[rankDate+"_ranking_2"]
		r3 := list[rankDate+"_ranking_3"]
		if i == dayDiff {
			dc.SetRGB(1, 0, 0)
		}
		dc.DrawString(r1, float64(x_offset+x_step*i+x_extra*i), float64(y_offset))
		dc.DrawString(r2, float64(x_offset+x_step*i+x_extra*i), float64(y_offset+y_step*1))
		dc.DrawString(r3, float64(x_offset+x_step*i+x_extra*i), float64(y_offset+y_step*2))
	}
	err = dc.SavePNG(savePath)
	return savePath, err
}
