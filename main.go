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
package main

import "sifrank/sched"

func main() {
	go sched.FetchNtpData()
	go sched.FetchPacketData()
	go sched.FetchRankData()

	select {}
}
