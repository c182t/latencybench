package bench

import (
	"time"
)

type MemoryStrideBenchmark struct {
	Options *BenchmarkOptions
}

func (msb *MemoryStrideBenchmark) Setup() error {
	return nil
}

func (msb *MemoryStrideBenchmark) RunOnce() (time.Duration, error) {
	src := make([]byte, msb.GetOptions().BlockSize*msb.GetOptions().Stride)
	sum := 0
	startTime := time.Now()
	for i := uint(0); i < uint(len(src)); i += msb.GetOptions().Stride {
		sum += int(src[i])
	}
	duration := time.Since(startTime)

	return duration, nil
}

func (msb *MemoryStrideBenchmark) Teardown() {}

func (msb *MemoryStrideBenchmark) Clone() Benchmark {
	clone := *msb
	return &clone
}

func (msb *MemoryStrideBenchmark) GetOptions() *BenchmarkOptions {
	return msb.Options
}
