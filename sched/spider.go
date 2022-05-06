//
// Copyright 2012 Google, Inc. All rights reserved.
// Copyright 2021-2022 YumeMichi. All rights reserved.
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
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"sifrank/config"
	"sifrank/db"
	"sifrank/xclog"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/examples/util"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

var ranking = ""

// Build a simple HTTP request parser using tcpassembly.StreamFactory and tcpassembly.Stream interfaces

// httpStreamFactory implements tcpassembly.StreamFactory
type httpStreamFactory struct{}

// httpStream will handle the actual decoding of http requests.
type httpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

func (h *httpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hstream.run() // Important... we must guarantee that data from the reader stream is read.

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &hstream.r
}

func (h *httpStream) run() {
	buf := bufio.NewReader(&h.r)
	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			// We must read until we see an EOF... very important!
			return
		} else if err != nil {
			// xclog.Warn("Error reading stream", h.net, h.transport, ":", err)
			continue
		} else {
			if req == nil {
				continue
			}
			// bodyBytes := tcpreader.DiscardBytesToEOF(req.Body)
			_ = req.Body.Close()
			// Filter URL - marathon
			if !strings.Contains(req.URL.String(), "eventPlayer") {
				continue
			}
			// Remove unused headers
			req.Header.Del("Debug")
			req.Header.Del("Connection")
			req.Header.Del("Content-Length")
			req.Header.Del("Accept-Encoding")
			req.Header.Del("Accept")
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 11; ONE E1001 Build/RQ1A.210105.003) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/87.0.4280.101 Mobile Safari/537.36")
			// Save headers to LevelDB
			headers := make(map[string]string)
			for k, v := range req.Header {
				headers[k] = v[0]
			}
			m, err := json.Marshal(headers)
			if err != nil {
				xclog.Warn("Failed to marshal headers into json: ", err.Error())
				continue
			}
			xclog.Info(string(m))
			if ranking == "" {
				continue
			}
			prefix := "request_header_"
			key := []byte(prefix + ranking)
			value := m
			err = db.LevelDb.Put(key, value)
			if err != nil {
				xclog.Warn(err.Error())
				continue
			}
			xclog.Info("Put ", string(key), " success!")
		}
	}
}

func FetchPacketData() {
	defer util.Run()()
	var handle *pcap.Handle
	var err error

	// Set up pcap packet capture
	if config.Conf.Fname != "" {
		xclog.Infof("Reading from pcap dump %s", config.Conf.Fname)
		handle, err = pcap.OpenOffline(config.Conf.Fname)
	} else {
		xclog.Infof("Starting capture on interface %s", config.Conf.Iface)
		handle, err = pcap.OpenLive(config.Conf.Iface, int32(config.Conf.Snaplen), true, pcap.BlockForever)
	}
	if err != nil {
		xclog.Error(err)
		return
	}

	if err := handle.SetBPFFilter(config.Conf.Filter); err != nil {
		xclog.Error(err)
		return
	}

	// Set up assembly
	streamFactory := &httpStreamFactory{}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	// Read in packets, pass to assembler.
	xclog.Info("Reading in packets")
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packets:
			// A nil packet indicates the end of a pcap file.
			if packet == nil {
				return
			}
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				xclog.Info("Unusable packet")
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			resp := string(tcp.Payload)
			kw := strings.Split(resp, "\n")
			for _, v := range kw {
				if !strings.Contains(v, "eventPlayer") || strings.Contains(v, "HTTP") {
					continue
				}
				if strings.Contains(v, "\"rank\":120") {
					ranking = "ranking_1"
				} else if strings.Contains(v, "\"rank\":700") {
					ranking = "ranking_2"
				} else if strings.Contains(v, "\"rank\":2300") {
					ranking = "ranking_3"
				} else {
					ranking = ""
				}
				if ranking == "" {
					continue
				}
				prefix := "request_data_"
				key := []byte(prefix + ranking)
				value := []byte(strings.TrimRight(v, "\r"))
				err = db.LevelDb.Put(key, value)
				if err != nil {
					xclog.Warn(err.Error())
					continue
				}
				xclog.Info("Put ", string(key), " success!")
			}
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		case <-ticker:
			// Every minute, flush connections that haven't seen activity in the past 2 minutes.
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}
