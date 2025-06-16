package bench

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type LoopbackTCPBenchmark struct {
	Options       *BenchmarkOptions
	serverContext context.Context
	cancelContext context.CancelFunc
	listener      net.Listener
	waitGroup     *sync.WaitGroup
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
			}
			return
		}
		conn.Write(buf[:n])
	}
}

func StartServer(protocol string, ip string, port string, ctx context.Context, wg *sync.WaitGroup) (net.Listener, error) {
	ln, err := net.Listen(protocol, fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		fmt.Printf("Unable to start server: %v", err)
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
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	ln, err := StartServer("tcp", "127.0.0.1", "9001", ctx, &wg)
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
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		return 0, fmt.Errorf("Unable to connect to 127.0.0.1: %v", err)
	}
	defer conn.Close()

	writeBuf := make([]byte, 4096)
	readBuf := make([]byte, 4096)

	startTime := time.Now()

	conn.Write(writeBuf)
	conn.Read(readBuf)

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
