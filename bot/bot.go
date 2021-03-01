package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io/ioutil"
	"log"
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

type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var ctx = context.Background()

func init() {
	rule := zero.SuffixRule("档线")
	zero.OnMessage(rule).SetBlock(true).SetPriority(10).
		Handle(func(ctx *zero.Ctx) {
			result, err := getData()
			if err != nil {
				log.Println(err)
				return
			}
			if len(result) != 3 {
				ctx.Send("【档线小助手】数据获取失败，请联系维护人员~")
				return
			}
			msg := fmt.Sprintf("档线小助手\n一档线: %d\n二档线: %d\n三档线: %d", result["ranking_1"], result["ranking_2"], result["ranking_3"])
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
		log.Println("No request header: ", err.Error())
		return map[string]int{}, err
	}
	for k, v := range result {
		requestData, err := rdb.HGet(ctx, "request_data", k).Result()
		if err == redis.Nil {
			log.Println("No request data for", k)
			continue
		} else if err != nil {
			log.Println("No request data: ", err.Error())
			continue
		}
		form := url.Values{"request_data": {requestData}}
		requestUrl := "http://prod.game1.ll.sdo.com/main.php/ranking/eventPlayer"
		req, err := http.NewRequest("POST", requestUrl, strings.NewReader(form.Encode()))
		if err != nil {
			log.Println("Send request error: ", err.Error())
			continue
		}
		headers := make(map[string]string)
		err = json.Unmarshal([]byte(v), &headers)
		if err != nil {
			log.Println("Unmarshal failed: ", err.Error())
			continue
		}
		for kk, vv := range headers {
			req.Header.Add(kk, vv)
		}

		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Send request failed: ", err.Error())
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)

		var res = &Response{}
		if err := json.Unmarshal(body, res); err != nil {
			log.Println("Unmarshal failed: ", err.Error())
			continue
		}
		items := res.ResponseData.Items
		itemLen := len(items)
		result := items[itemLen-1]
		ret[k] = result.Score

		_ = resp.Body.Close()
	}
	return ret, nil
}

