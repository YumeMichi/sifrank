package day

import (
	"sifrank/config"
	"sifrank/db"
	"sifrank/utils"
	"strconv"
	"time"

	"github.com/fogleman/gg"
)

var (
	x_offset = config.Conf.DqXOffset
	x_step   = config.Conf.DqXStep
	x_extra  = config.Conf.DqXExtra
	y_offset = config.Conf.DqYOffset
	y_step   = config.Conf.DqYStep
)

type DayRankData struct {
	Id       int    `db:"id"`
	Rank     string `db:"rank"`
	RankCode string `db:"rank_code"`
	Score    int    `db:"score"`
	DataDate string `db:"data_date"`
	DataTime string `db:"data_time"`
}

func GenDayRankPic() (string, error) {
	startDate, err := time.ParseInLocation("2006-01-02 15:04:05", config.Conf.StartTime, time.Local)
	fileName := startDate.Format("2006-01-02")
	if err != nil {
		return "", err
	}
	savePath := config.Conf.DqOutputDir + "/" + fileName + ".png"
	if utils.PathExists(config.Conf.DqOutputDir + "/" + fileName + ".png") {
		return savePath, nil
	}
	img, err := gg.LoadImage(config.Conf.DqBaseImage)
	if err != nil {
		return "", err
	}
	dc := gg.NewContextForImage(img)
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace(config.Conf.DqFontName, config.Conf.DqFontSize); err != nil {
		return "", err
	}
	// 档线大标题
	dc.DrawString(config.Conf.DqTitle, config.Conf.DqTitleXOffset, config.Conf.DqTitleYOffset)
	dc.SetRGB(1, 0, 0)
	if config.Conf.DqSubtitle != "" {
		dc.DrawString(config.Conf.DqSubtitle, config.Conf.DqSubtitleXOffset, config.Conf.DqSubtitleYOffset)
		dc.SetRGB(0, 0, 0)
	}
	// 档线小标题
	dc.DrawString("120", float64(x_offset), float64(y_offset))
	dc.DrawString("700", float64(x_offset), float64(y_offset+y_step*1))
	dc.DrawString("2300", float64(x_offset), float64(y_offset+y_step*2))
	// 档线数据
	for i := 1; i <= 12; i++ {
		timeDiff, _ := time.ParseDuration(strconv.Itoa(24*(i-1)) + "h")
		rankDate := startDate.Add(timeDiff).Format("2006-01-02")
		var dayRanks []DayRankData
		err := db.MysqlClient.Select(&dayRanks, "SELECT * FROM day_rank_data WHERE data_date = ? ORDER BY rank ASC", rankDate)
		if err != nil {
			panic(err)
		}
		if len(dayRanks) == 0 {
			break
		}
		if i == 12 {
			dc.SetRGB(1, 0, 0)
		}
		dc.DrawString(strconv.Itoa(dayRanks[0].Score), float64(x_offset+x_step*i+x_extra*i), float64(y_offset))
		dc.DrawString(strconv.Itoa(dayRanks[1].Score), float64(x_offset+x_step*i+x_extra*i), float64(y_offset+y_step*1))
		dc.DrawString(strconv.Itoa(dayRanks[2].Score), float64(x_offset+x_step*i+x_extra*i), float64(y_offset+y_step*2))
	}
	err = dc.SavePNG(savePath)
	return savePath, err
}
