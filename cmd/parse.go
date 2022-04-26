package cmd

import (
	"sifrank/db"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func init() {
	if len(GetKeywordList().Array()) == 0 {
		logrus.Warn("No keyword data")
	}
}

func GetKeywordList() (res gjson.Result) {
	key := []byte("keyword_list")
	list, err := db.LevelDb.Get(key)
	if err != nil {
		return
	}
	if gjson.ValidBytes(list) {
		res = gjson.ParseBytes(list)
	}
	return
}

func AddKeyword(res gjson.Result) (err error) {
	return nil
}

func DeleteKeywords(keywords []string) (err error) {
	return nil
}

func UpdateKeyword(keyword, content string) (err error) {
	return nil
}

func SearchKeywords(keywords []string) (res gjson.Result) {
	return
}

func IsKeywordExists(keyword string) bool {
	return false
}
