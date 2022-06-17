//
// Copyright (C) 2021 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package utils

import (
	"io/ioutil"
	"os"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func ReadAllText(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func WriteAllText(path, text string) {
	_ = ioutil.WriteFile(path, []byte(text), 0644)
}
