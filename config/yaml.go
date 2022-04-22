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
package config

import (
	"os"
	"sifrank/utils"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type YamlConfigs struct {
	Iface             string   `yaml:"iface"`
	Fname             string   `yaml:"fname"`
	Snaplen           int      `yaml:"snaplen"`
	Filter            string   `yaml:"filter"`
	NickName          []string `yaml:"nickname"`
	SuperUsers        []int64  `yaml:"super_users"`
	AdminUser         string   `yaml:"admin_user"`
	CqhttpHost        string   `yaml:"cqhttp_host"`
	CqhttpPort        string   `yaml:"cqhttp_port"`
	AccessToken       string   `yaml:"access_token"`
	RedisHost         string   `yaml:"redis_host"`
	RedisPort         string   `yaml:"redis_port"`
	RedisPassword     string   `yaml:"redis_password"`
	RedisDb           int      `yaml:"redis_db"`
	MysqlHost         string   `yaml:"mysql_host"`
	MysqlPort         string   `yaml:"mysql_port"`
	MysqlUser         string   `yaml:"mysql_user"`
	MysqlPassword     string   `yaml:"mysql_password"`
	MysqlDb           string   `yaml:"mysql_db"`
	LevelDbPath       string   `yaml:"leveldb_path"`
	AppName           string   `yaml:"app_name"`
	EventName         string   `yaml:"event_name"`
	StartTime         string   `yaml:"start_time"`
	EndTime           string   `yaml:"end_time"`
	DqXOffset         int      `yaml:"dq_x_offset"`
	DqXStep           int      `yaml:"dq_x_step"`
	DqXExtra          int      `yaml:"dq_x_extra"`
	DqYOffset         int      `yaml:"dq_y_offset"`
	DqYStep           int      `yaml:"dq_y_step"`
	DqTitle           string   `yaml:"dq_title"`
	DqSubtitle        string   `yaml:"dq_subtitle"`
	DqFontName        string   `yaml:"dq_font_name"`
	DqFontSize        float64  `yaml:"dq_font_size"`
	DqTitleXOffset    float64  `yaml:"dq_title_x_offset"`
	DqTitleYOffset    float64  `yaml:"dq_title_y_offset"`
	DqSubtitleXOffset float64  `yaml:"dq_subtitle_x_offset"`
	DqSubtitleYOffset float64  `yaml:"dq_subtitle_y_offset"`
	DqBaseImage       string   `yaml:"dq_base_image"`
	DqOutputDir       string   `yaml:"dq_output_dir"`
}

func DefaultConfigs() *YamlConfigs {
	return &YamlConfigs{
		Iface:             "enp8s0",
		Fname:             "",
		Snaplen:           1600,
		Filter:            "tcp and port 80",
		NickName:          []string{"YumeMichi"},
		SuperUsers:        []int64{785569962, 1157490807},
		AdminUser:         "1157490807",
		CqhttpHost:        "127.0.0.1",
		CqhttpPort:        "6700",
		AccessToken:       "",
		RedisHost:         "127.0.0.1",
		RedisPort:         "6379",
		RedisPassword:     "",
		RedisDb:           0,
		MysqlHost:         "127.0.0.1",
		MysqlPort:         "3306",
		MysqlUser:         "root",
		MysqlPassword:     "",
		MysqlDb:           "sifrank",
		LevelDbPath:       "./sifrank.db",
		AppName:           "LoveLive! 国服档线小助手",
		EventName:         "",
		StartTime:         "",
		EndTime:           "",
		DqXOffset:         62,
		DqXStep:           160,
		DqXExtra:          14,
		DqYOffset:         178,
		DqYStep:           46,
		DqTitle:           "活动标题",
		DqSubtitle:        "活动子标题",
		DqFontName:        "./simsun.ttc",
		DqFontSize:        32,
		DqTitleXOffset:    400,
		DqTitleYOffset:    60,
		DqSubtitleXOffset: 1800,
		DqSubtitleYOffset: 60,
		DqBaseImage:       "./base.jpg",
		DqOutputDir:       "./temp",
	}
}

func Load(p string) *YamlConfigs {
	if !utils.PathExists(p) {
		_ = DefaultConfigs().Save(p)
	}
	c := YamlConfigs{}
	err := yaml.Unmarshal([]byte(utils.ReadAllText(p)), &c)
	if err != nil {
		logrus.Error("[LLSIF] 尝试加载配置文件失败: 读取文件失败！")
		logrus.Info("[LLSIF] 原配置文件已备份！")
		_ = os.Rename(p, p+".backup"+strconv.FormatInt(time.Now().Unix(), 10))
		_ = DefaultConfigs().Save(p)
	}
	c = YamlConfigs{}
	_ = yaml.Unmarshal([]byte(utils.ReadAllText(p)), &c)
	logrus.Info("[LLSIF] 配置加载完毕！")
	return &c
}

func (c *YamlConfigs) Save(p string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		logrus.Error("[LLSIF] 写入新的配置文件失败！")
		return err
	}
	utils.WriteAllText(p, string(data))
	return nil
}
