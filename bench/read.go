package bench

import (
	"fmt"
	"io"
	"os"
	"time"
)

type ReadBenchmark struct {
	fd *os.File
}

func (rb *ReadBenchmark) Open() error {
	f, err := os.Open("/etc/hostname")
	if err != nil {
		return fmt.Errorf("BenchmarkRead() - Failed to open a file: %v", err)
	}
	rb.fd = f
	return nil
}

func (rb *ReadBenchmark) BenchmarkOnce() (time.Duration, error) {
	rb.fd.Seek(0, io.SeekStart)
	buf := make([]byte, 128)
	startTime := time.Now()
	_, err := rb.fd.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("BenchmarkRead() - Failed to read a file descriptor: %v; error: %v", rb.fd, err)
	}
	duration := time.Since(startTime)

	return duration, nil
}

func (rb *ReadBenchmark) Close() {
	if rb.fd != nil {
		rb.fd.Remo
		rb.fd.Close()
	}
}
