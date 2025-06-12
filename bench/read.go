package bench

import (
	"fmt"
	"os"
	"time"
)

func BenchmarkRead() (time.Duration, error) {
	f, err := os.Open("/etc/hostname")
	if err != nil {
		return 0, fmt.Errorf("BenchmarkRead() - Failed to open a file: %v", err)
	}
	defer f.Close()

	buf := make([]byte, 128)
	startTime := time.Now()
	_, err = f.Read(buf)
	duration := time.Since(startTime)

	return duration, nil
}
