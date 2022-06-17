//
// Copyright (C) 2022 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package db

import (
	"fmt"
	"sifrank/config"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Conf.RedisHost, config.Conf.RedisPort),
		Password: config.Conf.RedisPassword,
		DB:       config.Conf.RedisDb,
	})
}
