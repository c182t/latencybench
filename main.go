package main

import (
	"flag"
	"fmt"
	"latencybench/bench"
)

func RunBenchmark(b bench.Benchmark, n int, opts ...bench.BenchmarkOption) bench.BenchmarkAggregatedResult {
	cfg := bench.BenchmarkConfig{
		Parallelism: 1,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	var err error
	var benchmarkAggResult bench.BenchmarkAggregatedResult

	switch {
	case cfg.Parallelism == 1:
		benchmarkAggResult, err = bench.RunBenchmarkSerial(b, n)
	case cfg.Parallelism > 1:
		benchmarkAggResult, err = bench.RunBenchmarkParallel(b, n, cfg.Parallelism)
	default:
		fmt.Printf("Error occurred in RunBenchmark - incorrect parallelism value: %d\n", cfg.Parallelism)
		err = fmt.Errorf("invalid parallelism value: %d", cfg.Parallelism)
	}

	if err != nil {
		fmt.Printf("Error occurred in RunBenchmark: %v\n", err)
	}

	return benchmarkAggResult
}

func main() {
	benchmarkLabel := flag.String("benchmark", "read", "Benchmark to execute (e.g. read, write, sync)")
	iterations := flag.Int("iterations", 1000, "Number of benchmark iterations")
	parallelism := flag.Int("parallelism", 1, "Number of threads to run in parallel")
	blockSize := flag.Uint("block_size", 4*1024, "Block size (e.g. 4096)")

	flag.Parse()

	var benchmarkLabelMap = map[string]bench.Benchmark{
		"read":  &bench.ReadBenchmark{},
		"write": &bench.WriteBenchmark{},
		"sync":  &bench.SyncBenchmark{},
	}

	benchmark := benchmarkLabelMap[*benchmarkLabel]
	benchmarkResult := RunBenchmark(benchmark,
		*iterations,
		bench.WithParallelism(*parallelism),
		bench.WithBlockSize(*blockSize))

	fmt.Printf("%s benchmark = %+v\n", *benchmarkLabel, benchmarkResult)
}
