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
package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"sifrank/utils"
	"strconv"
	"time"
)

type YamlConfigs struct {
	Iface         string   `yaml:"iface"`
	Fname         string   `yaml:"fname"`
	Snaplen       int      `yaml:"snaplen"`
	Filter        string   `yaml:"filter"`
	NickName      []string `yaml:"nickname"`
	SuperUsers    []string `yaml:"super_users"`
	AdminUser     string   `yaml:"admin_user"`
	CqhttpHost    string   `yaml:"cqhttp_host"`
	CqhttpPort    string   `yaml:"cqhttp_port"`
	AccessToken   string   `yaml:"access_token"`
	RedisHost     string   `yaml:"redis_host"`
	RedisPort     string   `yaml:"redis_port"`
	RedisPassword string   `yaml:"redis_password"`
	RedisDb       int      `yaml:"redis_db"`
	Groups        []string `yaml:"groups"`
	AppName       string   `yaml:"app_name"`
	EventName     string   `yaml:"event_name"`
	EndTime       string   `yaml:"end_time"`
	IconUrl       string   `yaml:"icon_url"`
}

func DefaultConfigs() *YamlConfigs {
	return &YamlConfigs{
		Iface:         "enp8s0",
		Fname:         "",
		Snaplen:       1600,
		Filter:        "tcp and port 80",
		NickName:      []string{"YumeMichi"},
		SuperUsers:    []string{"785569962", "1157490807"},
		AdminUser:     "1157490807",
		CqhttpHost:    "127.0.0.1",
		CqhttpPort:    "6700",
		AccessToken:   "",
		RedisHost:     "127.0.0.1",
		RedisPort:     "6379",
		RedisPassword: "",
		RedisDb:       0,
		Groups:        []string{},
		AppName:       "LoveLive! 国服档线小助手",
		EventName:     "",
		EndTime:       "",
		IconUrl:       "https://c-ssl.duitang.com/uploads/item/201906/07/20190607235250_wtjcy.thumb.1000_0.jpg",
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

