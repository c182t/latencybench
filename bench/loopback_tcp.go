package bench

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

type LoopbackTCPBenchmark struct {
	Options   *BenchmarkOptions
	tcpServer TCPServer
}

func (ltb *LoopbackTCPBenchmark) Setup() error {
	tcpServer := TCPServer{}
	err := tcpServer.Start("tcp", "127.0.0.1", 0)
	if err != nil {
		return fmt.Errorf("Unable to start TCP server for LoopbackTCPBenchmark: %v", err)
	}
	ltb.tcpServer = tcpServer
	return nil
}

func (ltb *LoopbackTCPBenchmark) RunOnce() (time.Duration, error) {
	addr := fmt.Sprintf("%s:%d", ltb.tcpServer.ip, ltb.tcpServer.port)
	conn, err := net.Dial(ltb.tcpServer.protocol, addr)
	if err != nil {
		return 0, fmt.Errorf("unable to connect to %s on %s; %v", addr, ltb.tcpServer.protocol, err)
	}

	defer conn.Close()

	writeBuf := bytes.Repeat([]byte{'.'}, 4096)

	writeBuf[0] = '['
	writeBuf[len(writeBuf)-1] = ']'

	readBuf := make([]byte, 4096)

	startTime := time.Now()

	conn.Write(writeBuf)

	tcpConn := conn.(*net.TCPConn)
	tcpConn.CloseWrite()

	for {
		n, err := conn.Read(readBuf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error occured while read on client side (read %d bytes): %v", n, err)
			}
			break
		}
	}

	//fmt.Printf("readBuf=%s-%s", readBuf[:5], readBuf[len(readBuf)-5:len(readBuf)])

	duration := time.Since(startTime)
	return duration, nil
}

func (ltb *LoopbackTCPBenchmark) Teardown() {
	err := ltb.tcpServer.Stop()
	if err != nil {
		fmt.Printf("error occurd at Teardown: %v", err)
	}
}

func (ltb *LoopbackTCPBenchmark) Clone() Benchmark {
	clone := *ltb
	return &clone
}

func (ltb *LoopbackTCPBenchmark) GetOptions() *BenchmarkOptions {
	return ltb.Options
}
