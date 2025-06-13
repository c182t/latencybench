package bench

import (
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

type WriteBenchmark struct {
	fd *os.File
}

func (wb *WriteBenchmark) Open() error {
	randBytes := make([]byte, 16)
	randString := hex.EncodeToString(randBytes)

	f, err := os.OpenFile(fmt.Sprintf("/tmp/latencybench.write.%s", randString), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("WriteBenchmark() - Failed to open a file for write: %v", err)
	}
	wb.fd = f
	return nil
}

func (wb *WriteBenchmark) BenchmarkOnce() (time.Duration, error) {
	return 0, nil
}

func (wb *WriteBenchmark) Close() {
	if wb.fd != nil {
		wb.fd.Close()
	}
}
