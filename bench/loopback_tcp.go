package bench

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type LoopbackTCPBenchmark struct {
	Options       *BenchmarkOptions
	protocol      string
	ip            string
	port          string
	serverContext context.Context
	cancelContext context.CancelFunc
	listener      net.Listener
	waitGroup     *sync.WaitGroup
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1028)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
			}
			break
		}
		conn.Write(buf[:n])
	}
}

func StartServer(protocol string, ip string, port string, ctx context.Context, wg *sync.WaitGroup) (net.Listener, error) {
	addr := fmt.Sprintf("%s:%s", ip, port)

	ln, err := net.Listen(protocol, addr)
	if err != nil {
		fmt.Printf("Unable to start server [protocol=%s; addr=%s]: %v", protocol, addr, err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					fmt.Println("Server Accept error: ", err)
					continue
				}
			}
			go handleConnection(conn)
		}
	}()

	return ln, nil
}

func (ltb *LoopbackTCPBenchmark) Setup() error {
	ltb.protocol = "tcp"
	ltb.ip = "127.0.0.1"
	ltb.port = "9001"

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	ln, err := StartServer(ltb.protocol,
		ltb.ip, string(ltb.port), ctx, &wg)

	if err != nil {
		return fmt.Errorf("Unable to start TCP server for LoopbackTCPBenchmark: %v", err)
	}

	ltb.serverContext = ctx
	ltb.cancelContext = cancel
	ltb.listener = ln
	ltb.waitGroup = &wg

	return nil
}

func (ltb *LoopbackTCPBenchmark) RunOnce() (time.Duration, error) {
	addr := fmt.Sprintf("%s:%s", ltb.ip, ltb.port)
	conn, err := net.Dial(ltb.protocol, addr)
	if err != nil {
		return 0, fmt.Errorf("Unable to connect to %s on %s; %v", addr, ltb.protocol, err)
	}
	defer conn.Close()

	writeBuf := bytes.Repeat([]byte{'.'}, 4096)

	writeBuf[0] = '['
	writeBuf[len(writeBuf)-1] = ']'

	readBuf := make([]byte, 4096)

	startTime := time.Now()

	conn.Write(writeBuf)
	conn.Read(readBuf)
	//fmt.Printf("readBuf=%s-%s", readBuf[:5], readBuf[len(readBuf)-5:len(readBuf)])

	duration := time.Since(startTime)
	return duration, nil
}

func (ltb *LoopbackTCPBenchmark) Teardown() {
	if ltb.cancelContext != nil {
		ltb.cancelContext()

	}

	if ltb.listener != nil {
		ltb.listener.Close()
	}

	ltb.waitGroup.Wait()
}

func (ltb *LoopbackTCPBenchmark) Clone() Benchmark {
	clone := *ltb
	return &clone
}

func (ltb *LoopbackTCPBenchmark) GetOptions() *BenchmarkOptions {
	return ltb.Options
}
