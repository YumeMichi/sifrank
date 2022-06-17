//
// Copyright (C) 2021-2022 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package config

import (
	"os"
	"sifrank/utils"
	"sifrank/xclog"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type AppConfigs struct {
	AppName         string         `yaml:"app_name"`
	Log             LogConfigs     `yaml:"log"`
	Pcap            PcapConfigs    `yaml:"pcap"`
	Bot             BotConfigs     `yaml:"bot"`
	Redis           RedisConfigs   `yaml:"redis"`
	Mysql           MysqlConfigs   `yaml:"mysql"`
	LevelDb         LevelDbConfigs `yaml:"leveldb"`
	Event           EventConfigs   `yaml:"event"`
	Dq              DqConfigs      `yaml:"dq"`
	EnableMigration bool           `yaml:"enable_migration"`
}

type LogConfigs struct {
	LogDir   string `yaml:"log_dir"`
	LogLevel int    `yaml:"log_level"`
	LogSave  bool   `yaml:"log_save"`
}

type PcapConfigs struct {
	Iface   string `yaml:"iface"`
	Fname   string `yaml:"fname"`
	Snaplen int    `yaml:"snaplen"`
	Filter  string `yaml:"filter"`
}

type BotConfigs struct {
	NickName    []string `yaml:"nickname"`
	SuperUsers  []int64  `yaml:"super_users"`
	AdminUser   string   `yaml:"admin_user"`
	CqhttpHost  string   `yaml:"cqhttp_host"`
	CqhttpPort  string   `yaml:"cqhttp_port"`
	AccessToken string   `yaml:"access_token"`
}

type RedisConfigs struct {
	RedisHost     string `yaml:"redis_host"`
	RedisPort     string `yaml:"redis_port"`
	RedisPassword string `yaml:"redis_password"`
	RedisDb       int    `yaml:"redis_db"`
}

type MysqlConfigs struct {
	MysqlHost     string `yaml:"mysql_host"`
	MysqlPort     string `yaml:"mysql_port"`
	MysqlUser     string `yaml:"mysql_user"`
	MysqlPassword string `yaml:"mysql_password"`
	MysqlDb       string `yaml:"mysql_db"`
}

type LevelDbConfigs struct {
	DataPath string `yaml:"data_path"`
}

type EventConfigs struct {
	EventName string `yaml:"event_name"`
	StartTime string `yaml:"start_time"`
	EndTime   string `yaml:"end_time"`
}

type DqConfigs struct {
	DqXOffset         int     `yaml:"dq_x_offset"`
	DqXStep           int     `yaml:"dq_x_step"`
	DqXExtra          int     `yaml:"dq_x_extra"`
	DqYOffset         int     `yaml:"dq_y_offset"`
	DqYStep           int     `yaml:"dq_y_step"`
	DqTitle           string  `yaml:"dq_title"`
	DqSubtitle        string  `yaml:"dq_subtitle"`
	DqFontName        string  `yaml:"dq_font_name"`
	DqFontSize        float64 `yaml:"dq_font_size"`
	DqTitleXOffset    float64 `yaml:"dq_title_x_offset"`
	DqTitleYOffset    float64 `yaml:"dq_title_y_offset"`
	DqSubtitleXOffset float64 `yaml:"dq_subtitle_x_offset"`
	DqSubtitleYOffset float64 `yaml:"dq_subtitle_y_offset"`
	DqBaseImage       string  `yaml:"dq_base_image"`
	DqOutputDir       string  `yaml:"dq_output_dir"`
}

func DefaultConfigs() *AppConfigs {
	return &AppConfigs{
		AppName: "LoveLive! 国服档线小助手",
		Log: LogConfigs{
			LogDir:   "logs",
			LogLevel: 5,
			LogSave:  true,
		},
		Pcap: PcapConfigs{
			Iface:   "enp8s0",
			Fname:   "",
			Snaplen: 1600,
			Filter:  "tcp and port 80",
		},
		Bot: BotConfigs{
			NickName:    []string{"YumeMichi"},
			SuperUsers:  []int64{785569962, 1157490807},
			AdminUser:   "1157490807",
			CqhttpHost:  "127.0.0.1",
			CqhttpPort:  "6700",
			AccessToken: "",
		},
		Redis: RedisConfigs{
			RedisHost:     "127.0.0.1",
			RedisPort:     "6379",
			RedisPassword: "",
			RedisDb:       0,
		},
		Mysql: MysqlConfigs{
			MysqlHost:     "127.0.0.1",
			MysqlPort:     "3306",
			MysqlUser:     "root",
			MysqlPassword: "",
			MysqlDb:       "sifrank",
		},
		LevelDb: LevelDbConfigs{
			DataPath: "./sifrank.db",
		},
		Event: EventConfigs{
			EventName: "",
			StartTime: "",
			EndTime:   "",
		},
		Dq: DqConfigs{
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
		},
		EnableMigration: false,
	}
}

func Load(p string) *AppConfigs {
	if !utils.PathExists(p) {
		_ = DefaultConfigs().Save(p)
	}
	c := AppConfigs{}
	err := yaml.Unmarshal([]byte(utils.ReadAllText(p)), &c)
	if err != nil {
		xclog.Error("[LLSIF] 尝试加载配置文件失败: 读取文件失败！")
		xclog.Info("[LLSIF] 原配置文件已备份！")
		_ = os.Rename(p, p+".backup"+strconv.FormatInt(time.Now().Unix(), 10))
		_ = DefaultConfigs().Save(p)
	}
	c = AppConfigs{}
	_ = yaml.Unmarshal([]byte(utils.ReadAllText(p)), &c)
	xclog.Info("[LLSIF] 配置加载完毕！")
	return &c
}

func (c *AppConfigs) Save(p string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		xclog.Error("[LLSIF] 写入新的配置文件失败！")
		return err
	}
	utils.WriteAllText(p, string(data))
	return nil
}
