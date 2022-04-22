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
