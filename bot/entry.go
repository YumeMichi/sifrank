package bot

import (
	"fmt"
	"sifrank/db"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type KeywordEntry struct {
	Id      int64  `db:"id"`
	Keyword string `db:"keyword"`
	Content string `db:"content"`
	Sender  string `db:"sender"`
	Group   string `db:"group"`
	Dt      string `db:"dt"`
}

var (
	keywords []string
)

func InitKeywordEntry() {
	err := loadKeywords()
	if err != nil {
		return
	}

	// 添加词条
	entryRule := zero.PrefixRule("添加词条_")
	zero.OnMessage(entryRule, zero.OnlyToMe).SetBlock(true).SetPriority(10).
		Handle(func(context *zero.Ctx) {
			cq := message.At(context.Event.Sender.ID)
			msg := context.MessageString()
			msgList := strings.SplitN(msg, "_", 3)
			fmt.Println(msgList)
			if len(msgList) != 3 {
				context.Send(cq.String() + "格式有误！")
				return
			}
			key := msgList[1]
			value := msgList[2]
			sender := context.Event.Sender.ID
			group := context.Event.GroupID
			fmt.Println(key, value, sender, group)

			var entry []KeywordEntry
			err := db.MysqlClient.Select(&entry, "SELECT * FROM `keyword_entry` WHERE `keyword` = ? AND `group` = ?", key, group)
			if err != nil {
				logrus.Warn(err.Error())
				return
			}
			if len(entry) > 0 {
				context.Send(cq.String() + "该词条已存在！")
				return
			}
			dt := time.Now().Format("2006-01-02 15:04:05")
			ret, err := db.MysqlClient.Exec("INSERT INTO `keyword_entry` (`keyword`, `content`, `sender`, `group`, `dt`) VALUES (?, ?, ?, ?, ?)", key, value, sender, group, dt)
			if err != nil {
				logrus.Warn(err.Error())
				return
			}
			id, _ := ret.LastInsertId()
			logrus.Debug("Insert successfully. Id: ", id)
			loadKeywords()
			context.Send(cq.String() + "词条 " + key + " 添加成功！")
		})

	// 查找词条
	searchRule := zero.PrefixRule("查询词条")
	zero.OnMessage(searchRule, zero.OnlyToMe).SetBlock(true).SetPriority(10).
		Handle(func(context *zero.Ctx) {
			msg := context.MessageString()
			msgList := strings.Split(msg, "查询词条")
			key := strings.Trim(msgList[1], " ")
			group := context.Event.GroupID
			var entry []KeywordEntry
			err := db.MysqlClient.Select(&entry, "SELECT * FROM `keyword_entry` WHERE `keyword` LIKE ? AND `group` = ?", "%"+key+"%", group)
			if err != nil {
				logrus.Warn(err.Error())
				return
			}
			fmt.Println(entry)
		})
}

func loadKeywords() error {
	keywords = nil
	var entry []KeywordEntry
	err := db.MysqlClient.Select(&entry, "SELECT `keyword` FROM `keyword_entry`")
	if err != nil {
		return err
	}
	for _, v := range entry {
		keywords = append(keywords, v.Keyword)
	}
	logrus.Debug(keywords)
	reloadKeywordsSearch()
	return nil
}

func reloadKeywordsSearch() {
	// 完全匹配词条
	keywordRule := zero.FullMatchRule(keywords...)
	zero.OnMessage(keywordRule, zero.OnlyToMe).SetBlock(true).SetPriority(10).
		Handle(func(context *zero.Ctx) {
			key := context.MessageString()
			group := context.Event.GroupID
			var entry []KeywordEntry
			err := db.MysqlClient.Select(&entry, "SELECT * FROM `keyword_entry` WHERE `keyword` = ? AND `group` = ?", key, group)
			if err != nil {
				logrus.Warn(err.Error())
				return
			}
			cq := message.At(context.Event.Sender.ID)
			context.Send(cq.String() + "\n" + entry[0].Content)
		})
}
