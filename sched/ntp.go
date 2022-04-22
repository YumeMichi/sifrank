//
// Copyright 2022 YumeMichi. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// This binary provides sample code for using the gopacket TCP assembler and TCP
// stream reader.  It reads packets off the wire and reconstructs HTTP requests
// it sees, logging them.
//
package sched

import (
	"bytes"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
)

func FetchNtpData() {
	ticker := time.NewTicker(time.Hour * 6)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cmd := exec.Command("/usr/bin/ntpdate", "ntp.aliyun.com")
			cmd.Stderr = &bytes.Buffer{}
			cmd.Stdout = &bytes.Buffer{}
			err := cmd.Run()
			if err != nil {
				logrus.Warn(err.Error())
				logrus.Warn(cmd.Stderr.(*bytes.Buffer).String())
			} else {
				logrus.Info(cmd.Stdout.(*bytes.Buffer).String())
			}
		}
	}
}
