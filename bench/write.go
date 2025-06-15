package bench

import (
	"fmt"
	"os"
	"time"
)

type WriteBenchmark struct {
	fd       *os.File
	filePath string
}

func (wb *WriteBenchmark) Open() error {
	randString, err := RandomHexString(16)
	if err != nil {
		return fmt.Errorf("WriteBenchmark() - Failed to generate random file name: %v", err)
	}

	wb.filePath = fmt.Sprintf("/tmp/latencybench.write.%s", randString)

	f, err := os.OpenFile(wb.filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("WriteBenchmark() - Failed to open a file [%s] for write: %v", wb.filePath, err)
	}
	wb.fd = f
	return nil
}

func (wb *WriteBenchmark) BenchmarkOnce() (time.Duration, error) {
	startTime := time.Now()
	_, err := wb.fd.WriteString("0123456789012345")

	if err != nil {
		return 0, fmt.Errorf("WriteBenchmark() - Failed to write to a file: %v", err)
	}

	duration := time.Since(startTime)
	return duration, nil
}

func (wb *WriteBenchmark) Close() {
	if wb.fd != nil {
		wb.fd.Close()
		err := os.Remove(wb.filePath)
		if err != nil {
			fmt.Printf("WriteBenchmark() - Failed to remove file [%s]: %v", wb.filePath, err)
		}
	}
}
