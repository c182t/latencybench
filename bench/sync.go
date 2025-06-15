package bench

import (
	"fmt"
	"os"
	"time"
)

type SyncBenchmark struct {
	fd        *os.File
	filePath  string
	blockSize int
}

func (sb *SyncBenchmark) Setup() error {
	if sb.blockSize <= 0 {
		sb.blockSize = 4 * 1024
	}

	randFileName, err := RandomHexString(16)
	if err != nil {
		return fmt.Errorf("SyncBenchmark - unable to generate random file name: %v", err)
	}

	sb.filePath = fmt.Sprintf("/tmp/latencybench.sync.%v", randFileName)

	f, err := os.OpenFile(sb.filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("SyncBenchmark - unable to open file [%s] for writing: %v", sb.filePath, err)
	}

	sb.fd = f
	return nil
}

func (sb *SyncBenchmark) RunOnce() (time.Duration, error) {
	buf := make([]byte, sb.blockSize)

	_, err := sb.fd.Write(buf)

	if err != nil {
		return 0, fmt.Errorf("SyncBenchmark - Failed to write to a file [%s]: %v", sb.filePath, err)
	}
	startTime := time.Now()
	err = sb.fd.Sync()
	if err != nil {
		return 0, fmt.Errorf("SyncBenchmark - failed to sync to file [%s]: %v", sb.filePath, err)
	}
	duration := time.Since(startTime)
	return duration, nil
}

func (sb *SyncBenchmark) Teardown() {
	if sb.fd != nil {
		sb.fd.Close()
		err := os.Remove(sb.filePath)
		if err != nil {
			fmt.Printf("SyncBenchmark - Failed to remove file [%s]: %v", sb.filePath, err)
		}
	}
}

func (sb *SyncBenchmark) Clone() Benchmark {
	clone := *sb
	return &clone
}
