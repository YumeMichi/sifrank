//
// Copyright (C) 2021 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package config

var Conf = &AppConfigs{}

func init() {
	Conf = Load("./config.yml")
}
