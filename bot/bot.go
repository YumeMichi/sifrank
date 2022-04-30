//
// Copyright 2021-2022 YumeMichi. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// This binary provides sample code for using the gopacket TCP assembler and TCP
// stream reader.  It reads packets off the wire and reconstructs HTTP requests
// it sees, logging them.
//
package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sifrank/cmd"
	"sifrank/config"
	"sifrank/consts"
	"sifrank/db"
	"sifrank/tools"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type Response struct {
	ResponseData ResponseData `json:"response_data"`
	ReleaseInfo  interface{}  `json:"release_info"`
	StatusCode   int          `json:"status_code"`
}

type ResponseData struct {
	TotalCount      int         `json:"total_cnt"`
	Page            int         `json:"page"`
	Rank            interface{} `json:"rank"`
	Items           []ItemData  `json:"items"`
	ServerTimestamp int         `json:"server_timestamp"`
	PresentCount    int         `json:"present_cnt"`
}

type ItemData struct {
	Rank           int         `json:"rank"`
	Score          int         `json:"score"`
	UserData       interface{} `json:"user_data"`
	CenterUnitInfo interface{} `json:"center_unit_info"`
	SettingAwardId int         `json:"setting_award_id"`
}

var (
	limiter = rate.NewManager[int64](time.Second*5, 1)
)

func init() {
	zero.Run(zero.Config{
		NickName:      config.Conf.NickName,
		CommandPrefix: "/",
		SuperUsers:    config.Conf.SuperUsers,
		Driver: []zero.Driver{
			driver.NewWebSocketClient(fmt.Sprintf("ws://%s:%s", config.Conf.CqhttpHost, config.Conf.CqhttpPort), config.Conf.AccessToken),
		},
	})

	engine := zero.New()
	engine.UsePreHandler(func(ctx *zero.Ctx) bool {
		if !limiter.Load(ctx.Event.GroupID).Acquire() {
			ctx.Send("查询过于频繁！")
			return false
		}
		return true
	})

	rankRule := zero.FullMatchRule("档线", "dx")
	engine.OnMessage(rankRule).SetBlock(true).SetPriority(10).
		Handle(func(ctx *zero.Ctx) {
			now := time.Now()
			ed, err := time.ParseInLocation("2006-01-02 15:04:05", config.Conf.EndTime, time.Local)
			eds := ed.Format("20060102")
			if err != nil {
				logrus.Warn(err.Error())
				return
			}
			if now.After(ed) {
				list := db.LevelDb.ListPrefix([]byte(eds))
				r1 := list[eds+"_ranking_1"]
				r2 := list[eds+"_ranking_2"]
				r3 := list[eds+"_ranking_3"]
				// msg := fmt.Sprintf("【%s】\n当前活动: %s\n剩余时间: 已结束\n一档线点数: %s\n二档线点数: %s\n三档线点数: %s\n========================\n回复 dq/当期档线/本期档线 可查看每日档线数据", config.Conf.AppName, config.Conf.EventName, r1, r2, r3)
				msg := fmt.Sprintf("一档点数: %s\n二档点数: %s\n三档点数: %s\n剩余时间: 已结束", r1, r2, r3)
				ctx.Send(message.Text(msg))
				return
			}
			result, err := GetData()
			if err != nil || len(result) != 3 {
				logrus.Warn(err)
				dir, _ := os.Getwd()
				ctx.Send("【" + config.Conf.AppName + "】\n数据获取失败，请联系维护人员~\n[CQ:image,file=file:///" + filepath.ToSlash(filepath.Join(dir, "assets/images/emoji/fuck.jpg")) + "][CQ:at,qq=" + config.Conf.AdminUser + "]")
				return
			}
			// msg := fmt.Sprintf("【%s】\n当前活动: %s\n剩余时间: %s\n一档线点数: %s\n二档线点数: %s\n三档线点数: %s\n========================\n回复 dq/当期档线/本期档线 可查看每日档线数据", config.Conf.AppName, config.Conf.EventName, GetETA(), result["ranking_1"], result["ranking_2"], result["ranking_3"])
			msg := fmt.Sprintf("一档点数: %s\n二档点数: %s\n三档点数: %s\n剩余时间: %s", result["ranking_1"], result["ranking_2"], result["ranking_3"], GetETA())
			ctx.Send(message.Text(msg))
		})

	dayRankRule := zero.PrefixRule("当期", "当期档线", "本期档线", "dq")
	engine.OnMessage(dayRankRule).SetBlock(true).SetPriority(1).
		Handle(func(ctx *zero.Ctx) {
			savePath, err := cmd.GenDayRankPic()
			if err != nil {
				logrus.Warn(err.Error())
				return
			}

			dir, _ := os.Getwd()
			ctx.Send("[CQ:image,file=file:///" + filepath.ToSlash(filepath.Join(dir, savePath)) + "]")
		})

	engine.OnCommand("migrate", zero.AdminPermission).SetBlock(true).SetPriority(1).
		Handle(func(ctx *zero.Ctx) {
			tools.MigrateFromMySQLToLevelDB()
		})

	engine.OnCommand("list", zero.AdminPermission).SetBlock(true).SetPriority(1).
		Handle(func(ctx *zero.Ctx) {
			list := db.LevelDb.List()
			for k, v := range list {
				fmt.Println(string(k) + " - " + string(v))
			}
		})
}

func GetData() (map[string]string, error) {
	ret := make(map[string]string)
	for _, v := range consts.RankType {
		data_prefix := "request_data_"
		data_key := []byte(data_prefix + v)
		requestData, err := db.LevelDb.Get(data_key)
		if err != nil {
			logrus.Warn(err.Error())
			continue
		}
		form := url.Values{"request_data": {string(requestData)}}
		requestUrl := "http://prod.game1.ll.sdo.com/main.php/ranking/eventPlayer"
		req, err := http.NewRequest("POST", requestUrl, strings.NewReader(form.Encode()))
		if err != nil {
			logrus.Warn("Send request error: ", err.Error())
			continue
		}
		header_prefix := "request_header_"
		header_key := []byte(header_prefix + v)
		requestHeader, err := db.LevelDb.Get(header_key)
		if err != nil {
			logrus.Warn(err.Error())
			continue
		}
		headers := make(map[string]string)
		err = json.Unmarshal(requestHeader, &headers)
		if err != nil {
			logrus.Warn("Unmarshal failed: ", err.Error())
			continue
		}
		for kk, vv := range headers {
			req.Header.Add(kk, vv)
		}

		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			logrus.Warn("Send request failed: ", err.Error())
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)

		var res = &Response{}
		if err := json.Unmarshal(body, res); err != nil {
			logrus.Warn("Unmarshal failed: ", err.Error())
			continue
		}

		items := res.ResponseData.Items
		itemLen := len(items)
		if itemLen == 0 {
			ret[v] = "暂无数据"
		} else {
			result := items[itemLen-1]
			ret[v] = strconv.Itoa(result.Score)
		}

		_ = resp.Body.Close()

		time.Sleep(time.Millisecond * 300)
	}
	return ret, nil
}

func GetETA() string {
	now, err := ntp.Time("ntp.aliyun.com")
	if err != nil {
		logrus.Warn("NTP error, now using local time.")
		now = time.Now().Local()
	}
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", config.Conf.EndTime, time.Local)
	if now.After(end) {
		return "已结束"
	}
	hours := math.Floor(end.Sub(now).Hours())
	minutes := math.Floor(end.Sub(now).Minutes() - hours*60)
	if hours > 0 {
		return fmt.Sprintf("%d 小时 %d 分钟", int(hours), int(minutes))
	} else {
		return fmt.Sprintf("%d 分钟", int(minutes))
	}
}
