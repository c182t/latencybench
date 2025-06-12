package bench

import (
	"fmt"
	"os"
	"time"
)

type Benchmark interface {
	Open() error
	RunBenchmark() (time.Duration, error)
	Close()
}

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

func (rb *ReadBenchmark) Close() {
	if rb.fd != nil {
		rb.fd.Close()
	}
}

func (rb *ReadBenchmark) RunBenchmark() (time.Duration, error) {
	buf := make([]byte, 128)
	startTime := time.Now()
	_, err := rb.fd.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("BenchmarkRead() - Failed to read a file descriptor: %v; error: %v", rb.fd, err)
	}
	duration := time.Since(startTime)

	return duration, nil
}
