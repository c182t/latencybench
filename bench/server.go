package bench

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
)

type Server interface {
	Start(protocol string, ip string, port string) error
	Stop() error
}

type TCPServer struct {
	protocol string
	ip       string
	port     int
	listener net.Listener
	ctx      context.Context
	cancel   context.CancelFunc
	wg       *sync.WaitGroup
	status   string
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 2048)
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

func (svr *TCPServer) Start(protocol string, ip string, port int) error {
	ctx, cancel := context.WithCancel(context.Background())

	svr.ctx = ctx
	svr.cancel = cancel

	var wg sync.WaitGroup
	svr.wg = &wg

	addr := fmt.Sprintf("%s:%d", ip, port)

	ln, err := net.Listen(protocol, addr)
	if err != nil {
		fmt.Printf("Unable to start server [protocol=%s; addr=%s]: %v", protocol, addr, err)
	}

	svr.port = ln.Addr().(*net.TCPAddr).Port
	svr.ip = ip
	svr.protocol = protocol
	svr.listener = ln

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-svr.ctx.Done():
					return
				default:
					fmt.Println("Server: default Accept error: ", err)
					continue
				}
			}
			go handleConnection(conn)
		}
	}()

	svr.status = "started"
	return nil
}

func (svr *TCPServer) Stop() error {
	if svr.cancel != nil {
		svr.cancel()
	}

	if svr.listener != nil {
		svr.listener.Close()
	}

	svr.wg.Wait()
	svr.status = "stopped"

	return nil
}
