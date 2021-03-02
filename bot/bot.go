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

var ctx = context.Background()

func init() {
	rule := zero.SuffixRule("档线")
	zero.OnMessage(rule).SetBlock(true).SetPriority(10).
		Handle(func(ctx *zero.Ctx) {
			result, err := getData()
			if err != nil || len(result) != 3 {
				logrus.Warn(err.Error())
				ctx.Send("【LoveLive! 国服档线小助手】\n数据获取失败，请联系维护人员~")
				return
			}
			msg := fmt.Sprintf("【LoveLive! 国服档线小助手】\n当前活动: AZALEA 的前进之路!\n剩余时间: %s\n一档线积分: %d\n二档线积分: %d\n三档线积分: %d", getETA(), result["ranking_1"], result["ranking_2"], result["ranking_3"])
			ctx.Send(message.Text(msg))
		})
}

func getData() (map[string]int, error) {
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

func getETA() string {
	now := time.Now().Local()
	end, _ := time.Parse("2006-01-02 15:04:05", "2021-03-08 14:00:00")
	hours := math.Floor(end.Sub(now).Hours())
	minutes := math.Floor(end.Sub(now).Minutes() - hours*60)
	if hours > 0 {
		return fmt.Sprintf("%d 小时 %d 分钟", int(hours), int(minutes))
	} else {
		return fmt.Sprintf("%d 分钟", int(minutes))
	}
}

