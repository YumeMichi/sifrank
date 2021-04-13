//
// Copyright 2021 YumeMichi. All rights reserved.
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
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sifrank/config"
	"strconv"
	"strings"
	"time"
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

type CardInfo struct {
	App         string `json:"app"`
	Desc        string `json:"desc"`
	View        string `json:"view"`
	Ver         string `json:"ver"`
	Prompt      string `json:"prompt"`
	AppID       string `json:"appID"`
	SourceName  string `json:"sourceName"`
	ActionData  string `json:"actionData"`
	ActionDataA string `json:"actionData_A"`
	SourceUrl   string `json:"sourceUrl"`
	Meta        struct {
		Notification struct {
			AppInfo struct {
				AppName string `json:"appName"`
				AppType int    `json:"appType"`
				AppId   int    `json:"appid"`
				IconUrl string `json:"iconUrl"`
			} `json:"appInfo"`
			Data [4]struct {
				Title string `json:"title"`
				Value string `json:"value"`
			} `json:"data"`
			Title  string `json:"title"`
			Button [0]struct {
				Name   string `json:"name"`
				Action string `json:"action"`
			} `json:"button"`
			EmphasisKeyword string `json:"emphasis_keyword"`
		} `json:"notification"`
	} `json:"meta"`
	Text     string `json:"text"`
	SourceAd string `json:"sourceAd"`
}

var ctx = context.Background()
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Conf.RedisHost, config.Conf.RedisPort),
		Password: config.Conf.RedisPassword,
		DB:       config.Conf.RedisDb,
	})

	rankRule := zero.FullMatchRule("档线", "dx")
	zero.OnMessage(rankRule).SetBlock(true).SetPriority(10).
		Handle(func(context *zero.Ctx) {
			lock, _ := rdb.Get(ctx, "dx_lock").Result()
			if lock != "" {
				context.Send(message.Text("查询过于频繁！"))
				return
			}
			err := rdb.Set(ctx, "dx_lock", "1", time.Second*3).Err()
			if err != nil {
				logrus.Warn(err.Error())
				return
			}
			result, err := GetData()
			if err != nil || len(result) != 3 {
				logrus.Warn(err)
				dir, _ := os.Getwd()
				context.Send("【" + config.Conf.AppName + "】\n数据获取失败，请联系维护人员~\n[CQ:image,file=file:///" + filepath.ToSlash(filepath.Join(dir, "assets/images/emoji/fuck.jpg")) + "][CQ:at,qq=" + config.Conf.AdminUser + "]")
				return
			}
			msg := fmt.Sprintf("【%s】\n当前活动: %s\n剩余时间: %s\n一档线点数: %s\n二档线点数: %s\n三档线点数: %s", config.Conf.AppName, config.Conf.EventName, GetETA(), result["ranking_1"], result["ranking_2"], result["ranking_3"])
			context.Send(message.Text(msg))
		})

	cardRule := zero.FullMatchRule("查询档线（暂停使用）")
	zero.OnMessage(cardRule).SetBlock(true).SetPriority(1).
		Handle(func(ctx *zero.Ctx) {
			result, err := GetData()
			if err != nil || len(result) != 3 {
				logrus.Warn(err)
				dir, _ := os.Getwd()
				ctx.Send("【" + config.Conf.AppName + "】\n数据获取失败，请联系维护人员~\n[CQ:image,file=file:///" + filepath.ToSlash(filepath.Join(dir, "assets/images/emoji/fuck.jpg")) + "][CQ:at,qq=" + config.Conf.AdminUser + "]")
				return
			}
			var card = &CardInfo{}
			card.App = "com.tencent.miniapp"
			card.View = "notification"
			card.Ver = "0.0.0.1"
			card.Prompt = "[应用]"
			card.Meta.Notification.AppInfo.AppName = config.Conf.AppName
			card.Meta.Notification.AppInfo.AppType = 4
			card.Meta.Notification.AppInfo.AppId = 1109659848
			card.Meta.Notification.AppInfo.IconUrl = config.Conf.IconUrl
			card.Meta.Notification.Data[0].Title = "结束时间"
			card.Meta.Notification.Data[0].Value = GetETA()
			card.Meta.Notification.Data[1].Title = "一档线点数"
			card.Meta.Notification.Data[1].Value = result["ranking_1"]
			card.Meta.Notification.Data[2].Title = "二档线点数"
			card.Meta.Notification.Data[2].Value = result["ranking_2"]
			card.Meta.Notification.Data[3].Title = "三档线点数"
			card.Meta.Notification.Data[3].Value = result["ranking_3"]
			card.Meta.Notification.Title = config.Conf.EventName
			msg, err := json.Marshal(card)
			if err != nil {
				logrus.Warn("Marshal failed: ", err.Error())
				return
			}
			content := strings.ReplaceAll(string(msg), ",", "&#44;")
			content = strings.ReplaceAll(content, "[", "&#91;")
			content = strings.ReplaceAll(content, "]", "&#93;")
			ctx.Send("[CQ:json,data=" + content + "]")
		})
}

func GetData() (map[string]string, error) {
	ret := make(map[string]string)
	result, err := rdb.HGetAll(ctx, "request_header").Result()
	if err != nil {
		logrus.Warn("No request header: ", err.Error())
		return map[string]string{}, err
	}
	for k, v := range result {
		requestData, err := rdb.HGet(ctx, "request_data", k).Result()
		if err == redis.Nil {
			logrus.Warn("No request data for", k)
			continue
		} else if err != nil {
			logrus.Warn("No request data: ", err.Error())
			continue
		}
		form := url.Values{"request_data": {requestData}}
		requestUrl := "http://prod.game1.ll.sdo.com/main.php/ranking/eventPlayer"
		req, err := http.NewRequest("POST", requestUrl, strings.NewReader(form.Encode()))
		if err != nil {
			logrus.Warn("Send request error: ", err.Error())
			continue
		}
		headers := make(map[string]string)
		err = json.Unmarshal([]byte(v), &headers)
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
			ret[k] = "暂无数据"
		} else {
			result := items[itemLen-1]
			ret[k] = strconv.Itoa(result.Score)
		}

		_ = resp.Body.Close()

		time.Sleep(time.Millisecond * 300)
	}
	return ret, nil
}

func GetETA() string {
	now := time.Now().Local()
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

