//
// Copyright (C) 2021 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package config

var Conf = &YamlConfigs{}

func init() {
	Conf = Load("./config.yml")
}
