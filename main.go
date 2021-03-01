// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// This binary provides sample code for using the gopacket TCP assembler and TCP
// stream reader.  It reads packets off the wire and reconstructs HTTP requests
// it sees, logging them.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"github.com/go-redis/redis/v8"
	"github.com/google/gopacket"
	"github.com/google/gopacket/examples/util"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"io"
	"log"
	"net/http"
	_ "sifrank/bot"
	"strings"
	"time"
)

var iface = flag.String("i", "enp8s0", "Interface to get packets from")
var fname = flag.String("r", "", "Filename to read from, overrides -i")
var snaplen = flag.Int("s", 1600, "SnapLen for pcap packet capture")
var filter = flag.String("f", "tcp and port 80", "BPF filter for pcap")

var ctx = context.Background()
var ranking = ""

var rdb *redis.Client

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
			continue
			// log.Println("Error reading stream", h.net, h.transport, ":", err)
		} else {
			if req == nil {
				continue
			}
			//bodyBytes := tcpreader.DiscardBytesToEOF(req.Body)
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
			// Save headers to Redis
			headers := make(map[string]string)
			for k, v := range req.Header {
				headers[k] = v[0]
			}
			m, err := json.Marshal(headers)
			if err != nil {
				log.Println("Failed to marshal headers into json: ", err.Error())
				continue
			}
			log.Println(string(m))
			if ranking == "" {
				continue
			}
			requestHeader := map[string]interface{}{
				ranking: string(m),
			}
			err = rdb.HSet(ctx, "request_header", requestHeader).Err()
			if err != nil {
				log.Println("Failed to save json to Redis: ", err.Error())
				continue
			}
			log.Println("Success!")
		}
	}
}

func init() {
	logrus.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[YumeMichi][%time%][%lvl%]: %msg% \n",
	})
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	defer util.Run()()
	var handle *pcap.Handle
	var err error

	// Set up pcap packet capture
	if *fname != "" {
		log.Printf("Reading from pcap dump %q", *fname)
		handle, err = pcap.OpenOffline(*fname)
	} else {
		log.Printf("Starting capture on interface %q", *iface)
		handle, err = pcap.OpenLive(*iface, int32(*snaplen), true, pcap.BlockForever)
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := handle.SetBPFFilter(*filter); err != nil {
		log.Fatal(err)
	}

	// Set up assembly
	streamFactory := &httpStreamFactory{}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	// Connect to Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	// Init bot
	zero.Run(zero.Config{
		NickName:      []string{"YumeMichi"},
		CommandPrefix: "/",
		SuperUsers:    []string{"785569962", "1157490807"},
		Driver: []zero.Driver{
			driver.NewWebSocketClient("127.0.0.1", "6700", ""),
		},
	})

	log.Println("Reading in packets")
	// Read in packets, pass to assembler.
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
				log.Println("Unusable packet")
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			resp := string(tcp.Payload)
			kw := strings.Split(resp, "\n")
			for _, v := range kw {
				if !strings.Contains(v, "eventPlayer") || strings.Contains(v, "HTTP") {
					continue
				}
				log.Println(v)
				if strings.Contains(v, "\"rank\":\"120\"") {
					ranking = "ranking_1"
				} else if strings.Contains(v, "\"rank\":\"700\"") {
					ranking = "ranking_2"
				} else if strings.Contains(v, "\"rank\":\"2300\"") {
					ranking = "ranking_3"
				} else {
					ranking = ""
				}
				if ranking == "" {
					continue
				}
				requestData := map[string]interface{}{
					ranking: strings.TrimRight(v, "\r"),
				}
				err = rdb.HSet(ctx, "request_data", requestData).Err()
				if err != nil {
					log.Println("Failed to save json to Redis: ", err.Error())
					continue
				}
			}
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		case <-ticker:
			// Every minute, flush connections that haven't seen activity in the past 2 minutes.
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}

