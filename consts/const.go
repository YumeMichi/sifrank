//
// Copyright (C) 2022 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package consts

var (
	RankCode = make(map[string]int)
	RankType = []string{}
)

func init() {
	RankCode["ranking_1"] = 120
	RankCode["ranking_2"] = 700
	RankCode["ranking_3"] = 2300

	RankType = append(RankType, "ranking_1")
	RankType = append(RankType, "ranking_2")
	RankType = append(RankType, "ranking_3")
}
