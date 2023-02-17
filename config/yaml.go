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
	LevelDb         LevelDbConfigs `yaml:"leveldb"`
	Event           EventConfigs   `yaml:"event"`
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

type LevelDbConfigs struct {
	DataPath string `yaml:"data_path"`
}

type EventConfigs struct {
	StartTime string `yaml:"start_time"`
	EndTime   string `yaml:"end_time"`
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
		LevelDb: LevelDbConfigs{
			DataPath: "./sifrank.db",
		},
		Event: EventConfigs{
			StartTime: "",
			EndTime:   "",
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
