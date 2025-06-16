package bench

import (
	"time"
)

type RetIntBenchmark struct {
	Options *BenchmarkOptions
}

func (rib *RetIntBenchmark) Setup() error {
	return nil
}

func RetInt() int {
	return 5555555
}

func (rib *RetIntBenchmark) RunOnce() (time.Duration, error) {
	startTime := time.Now()
	someInt := RetInt()
	duration := time.Since(startTime)
	someInt = someInt + 0
	return duration, nil
}

func (rib *RetIntBenchmark) Teardown() {}

func (rib *RetIntBenchmark) Clone() Benchmark {
	clone := *rib
	return &clone
}

func (rib *RetIntBenchmark) GetOptions() *BenchmarkOptions {
	return rib.Options
}
