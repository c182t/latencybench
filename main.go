package main

import (
	"fmt"
	"latencybench/bench"
	"time"
)

type BenchmarkConfig struct {
	parallelism int
}

type BenchmarkOption func(*BenchmarkConfig)

func WithParallelism(parallelism int) BenchmarkOption {
	return func(bc *BenchmarkConfig) {
		bc.parallelism = parallelism
	}
}

func RunBenchmark(b bench.Benchmark, n int, opts ...BenchmarkOption) bench.BenchmarkResult {
	cfg := BenchmarkConfig{
		parallelism: 1,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	b.Open()
	defer b.Close()

	readBenchmarkFn := func() (time.Duration, error) {
		return b.BenchmarkOnce()
	}

	var err error
	var benchmarkResult bench.BenchmarkResult

	switch {
	case cfg.parallelism == 1:
		benchmarkResult, err = bench.RunBenchmark(readBenchmarkFn, n)
	case cfg.parallelism > 1:
		benchmarkResult, err = bench.RunBenchmarkParallel(readBenchmarkFn, n, cfg.parallelism)
	default:
		fmt.Printf("Error occurred in RunBenchmark - incorrect parallelism value: %d", cfg.parallelism)
		benchmarkResult = bench.BenchmarkResult{}
		err = fmt.Errorf("Ivalid parallelis value: %d", cfg.parallelism)
	}

	if err != nil {
		fmt.Printf("Error occurred in RunBenchmark: %v", err)
	}

	return benchmarkResult
}

func main() {
	fmt.Println("Syscall Latency Benchmarks: ")

	readBenchmarkResult := RunBenchmark(&bench.ReadBenchmark{}, 10000)
	fmt.Printf("read benchmark = %+v\n", readBenchmarkResult)

	writeBenchmarkResult := RunBenchmark(&bench.WriteBenchmark{}, 10000)
	fmt.Printf("write benchmark = %+v\n", writeBenchmarkResult)

	syncBenchmarkResult := RunBenchmark(&bench.SyncBenchmark{}, 10000)
	fmt.Printf("sync benchmark = %+v\n", syncBenchmarkResult)

	for _, i := range []int{2, 4, 8, 16, 32, 64} {
		syncBenchmarkResult := RunBenchmark(&bench.SyncBenchmark{}, 10000, WithParallelism(i))
		fmt.Printf("parallel sync [parallelism=%d] benchmark = %+v\n", i, syncBenchmarkResult)
	}

}
