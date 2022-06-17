//
// Copyright (C) 2022 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
//
package model

type DayRankData struct {
	Id       int    `db:"id"`
	Rank     string `db:"rank"`
	RankCode string `db:"rank_code"`
	Score    int    `db:"score"`
	DataDate string `db:"data_date"`
	DataTime string `db:"data_time"`
}

type EntryData struct {
	Id       int64  `json:"id"`
	Keyword  string `json:"keyword"`
	Content  string `json:"content"`
	Setter   int64  `json:"setter"`
	Group    int64  `json:"group"`
	DateTime string `json:"date_time"`
}
