// Copyright (c) 2019 Baidu, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// packetlog is packet-oriented log writer

// Using PacketConn interface in golang (net.PacketConn).
// The network net must be a packet-oriented network: udp, udp4, udp6, unixgram

package log4go

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

// packet connection
type PacketConn struct {
	conn       net.PacketConn
	remoteAddr net.Addr
}

var (
	ErrNetworkNotMatch = errors.New("network in following list:  udp, udp4, udp6, unixgram")
)

func resolveAddr(network string, remoteAddr string) (net.Addr, error) {
	var addr net.Addr
	var err error

	network = strings.ToLower(network)
	switch network {
	case "udp":
		fallthrough
	case "udp4":
		fallthrough
	case "udp6":
		addr, err = net.ResolveUDPAddr(network, remoteAddr)
	case "unixgram":
		addr, err = net.ResolveUnixAddr(network, remoteAddr)
	default:
		addr = nil
		err = ErrNetworkNotMatch
	}

	return addr, err
}

// newPacketConn creates packet connection
func newPacketConn(network string, remoteAddr string) (*PacketConn, error) {
	network = strings.ToLower(network)
	if network != "udp" && network != "udp4" &&
		network != "udp6" && network != "unixgram" {
		return nil, ErrNetworkNotMatch
	}
	// create golang net.PacketConn
	// Do not receive data, so local address is set to "".
	conn, err := net.ListenPacket(network, "")
	if err != nil {
		return nil, err
	}

	// resolve address
	address, err := resolveAddr(network, remoteAddr)
	if err != nil {
		return nil, err
	}

	// create PacketConn
	pc := &PacketConn{conn: conn, remoteAddr: address}

	return pc, nil
}

// Send sends data to log server
func (pc *PacketConn) Send(data []byte) error {
	_, err := pc.conn.WriteTo(data, pc.remoteAddr)
	return err
}

// PacketWriter sends output to a packet connection
type PacketWriter struct {
	LogCloser //for Elegant exit

	rec  chan *LogRecord
	conn *PacketConn
	name string
}

// Send sends data
func (w *PacketWriter) Send(data []byte) error {
	return w.conn.Send(data)
}

func (w *PacketWriter) LogWrite(rec *LogRecord) {
	if !LogWithBlocking {
		if len(w.rec) >= LogBufferLength {
			return
		}
	}

	w.rec <- rec
}

// Name gets writer name
func (w *PacketWriter) Name() string {
	return w.name
}

// QueueLen gets length of rec channel
func (w *PacketWriter) QueueLen() int {
	return len(w.rec)
}

func NewPacketWriter(name string, network string,
	remoteAddr string, format string) *PacketWriter {
	conn, err := newPacketConn(network, remoteAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewPacketWriter(%s, %s): %s\n",
			name, remoteAddr, err)
		return nil
	}

	w := &PacketWriter{
		rec:  make(chan *LogRecord, LogBufferLength),
		conn: conn,
		name: name,
	}

	//init LogCloser
	w.LogCloserInit()

	// add w to collection of all writers' info
	writersInfo = append(writersInfo, w)

	go func() {
		for {
			rec := <-w.rec

			if w.EndNotify(rec) {
				return
			}

			if rec.Binary != nil {
				w.Send(rec.Binary)
				putBuffer(rec.Binary) // Binary is allocated from buffer pool
			} else {
				msg := FormatLogRecord(format, rec)
				w.Send([]byte(msg))
			}
		}
	}()

	return w
}

// Close waits for dump all log and closes chan
func (w *PacketWriter) Close() {
	w.WaitForEnd(w.rec)
	close(w.rec)
}
