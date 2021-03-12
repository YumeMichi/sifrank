//
// Copyright 2021 YumeMichi. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// This binary provides sample code for using the gopacket TCP assembler and TCP
// stream reader.  It reads packets off the wire and reconstructs HTTP requests
// it sees, logging them.
package bot

import (
	"context"
	"encoding/json"
	"errors"
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

func init() {
	rankRule := zero.FullMatchRule("档线")
	zero.OnMessage(rankRule).SetBlock(true).SetPriority(10).
		Handle(func(ctx *zero.Ctx) {
			result, err := GetData()
			if err != nil || len(result) != 3 {
				logrus.Warn(err)
				dir, _ := os.Getwd()
				ctx.Send("【LoveLive! 国服档线小助手】\n数据获取失败，请联系维护人员~\n[CQ:image,file=file:///" + filepath.ToSlash(filepath.Join(dir, "assets/images/emoji/fuck.jpg")) + "][CQ:at,qq=1157490807]")
				return
			}
			msg := fmt.Sprintf("【LoveLive! 国服档线小助手】\n当前活动: AZALEA 的前进之路!\n剩余时间: %s\n一档线积分: %d\n二档线积分: %d\n三档线积分: %d", GetETA(), result["ranking_1"], result["ranking_2"], result["ranking_3"])
			ctx.Send(message.Text(msg))
		})

	cardRule := zero.FullMatchRule("查询档线")
	zero.OnMessage(cardRule).SetBlock(true).SetPriority(1).
		Handle(func(ctx *zero.Ctx) {
			result, err := GetData()
			if err != nil || len(result) != 3 {
				logrus.Warn(err)
				dir, _ := os.Getwd()
				ctx.Send("【LoveLive! 国服档线小助手】\n数据获取失败，请联系维护人员~\n[CQ:image,file=file:///" + filepath.ToSlash(filepath.Join(dir, "assets/images/emoji/fuck.jpg")) + "][CQ:at,qq=1157490807]")
				return
			}
			var card = &CardInfo{}
			card.App = "com.tencent.miniapp"
			card.View = "notification"
			card.Ver = "0.0.0.1"
			card.Prompt = "[应用]"
			card.Meta.Notification.AppInfo.AppName = "LoveLive! 国服档线小助手"
			card.Meta.Notification.AppInfo.AppType = 4
			card.Meta.Notification.AppInfo.AppId = 1109659848
			card.Meta.Notification.AppInfo.IconUrl = "https://c-ssl.duitang.com/uploads/item/201906/07/20190607235250_wtjcy.thumb.1000_0.jpg"
			card.Meta.Notification.Data[0].Title = "结束时间"
			card.Meta.Notification.Data[0].Value = GetETA()
			card.Meta.Notification.Data[1].Title = "一档线"
			card.Meta.Notification.Data[1].Value = strconv.Itoa(result["ranking_1"])
			card.Meta.Notification.Data[2].Title = "二档线"
			card.Meta.Notification.Data[2].Value = strconv.Itoa(result["ranking_2"])
			card.Meta.Notification.Data[3].Title = "三档线"
			card.Meta.Notification.Data[3].Value = strconv.Itoa(result["ranking_3"])
			card.Meta.Notification.Title = "AZALEA 的前进之路!"
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

func GetData() (map[string]int, error) {
	ret := make(map[string]int)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	result, err := rdb.HGetAll(ctx, "request_header").Result()
	if err != nil {
		logrus.Warn("No request header: ", err.Error())
		return map[string]int{}, err
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
			_ = resp.Body.Close()
			return map[string]int{}, errors.New(string(body))
		}
		result := items[itemLen-1]
		ret[k] = result.Score

		_ = resp.Body.Close()

		time.Sleep(time.Millisecond * 300)
	}
	return ret, nil
}

func GetETA() string {
	now := time.Now().Local()
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-03-08 14:00:00", time.Local)
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

