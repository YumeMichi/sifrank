//
// Copyright (C) 2022 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package modules

import (
	"bytes"
	"os/exec"
	"sifrank/xclog"
	"time"
)

func SyncNtpDate() {
	ticker := time.NewTicker(time.Hour * 6)
	defer ticker.Stop()

	for range ticker.C {
		cmd := exec.Command("/usr/bin/ntpdate", "ntp.aliyun.com")
		cmd.Stderr = &bytes.Buffer{}
		cmd.Stdout = &bytes.Buffer{}
		err := cmd.Run()
		if err != nil {
			xclog.Warn(err.Error())
			xclog.Warn(cmd.Stderr.(*bytes.Buffer).String())
		} else {
			xclog.Info(cmd.Stdout.(*bytes.Buffer).String())
		}
	}
}
