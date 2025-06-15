package bench

import (
	"time"
)

type MemoryCopyBenchmark struct {
	Options *BenchmarkOptions
}

func (mcb *MemoryCopyBenchmark) Setup() error {
	return nil
}

func (mcb *MemoryCopyBenchmark) RunOnce() (time.Duration, error) {
	src := make([]byte, mcb.GetOptions().BlockSize)
	dst := make([]byte, mcb.GetOptions().BlockSize)

	startTime := time.Now()
	copy(dst, src)
	duration := time.Since(startTime)

	return duration, nil
}

func (mcb *MemoryCopyBenchmark) Teardown() {}

func (mcb *MemoryCopyBenchmark) Clone() Benchmark {
	clone := *mcb
	return &clone
}

func (mcb *MemoryCopyBenchmark) GetOptions() *BenchmarkOptions {
	return mcb.Options
}
