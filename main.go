package main

import (
	"flag"
	"fmt"
	"latencybench/bench"
)

type BenchmarkConfig struct {
	benchmark   string
	iterations  int
	parallelism int
}

type BenchmarkOption func(*BenchmarkConfig)

func WithBenchmark(benchmark string) BenchmarkOption {
	return func(bc *BenchmarkConfig) {
		bc.benchmark = benchmark
	}
}

func WithIterations(iterations int) BenchmarkOption {
	return func(bc *BenchmarkConfig) {
		bc.iterations = iterations
	}
}

func WithParallelism(parallelism int) BenchmarkOption {
	return func(bc *BenchmarkConfig) {
		bc.parallelism = parallelism
	}
}

func RunBenchmark(b bench.Benchmark, n int, opts ...BenchmarkOption) bench.BenchmarkAggregatedResult {
	cfg := BenchmarkConfig{
		parallelism: 1,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	var err error
	var benchmarkAggResult bench.BenchmarkAggregatedResult

	switch {
	case cfg.parallelism == 1:
		benchmarkAggResult, err = bench.RunBenchmarkSerial(b, n)
	case cfg.parallelism > 1:
		benchmarkAggResult, err = bench.RunBenchmarkParallel(b, n, cfg.parallelism)
	default:
		fmt.Printf("Error occurred in RunBenchmark - incorrect parallelism value: %d\n", cfg.parallelism)
		err = fmt.Errorf("invalid parallelism value: %d", cfg.parallelism)
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

	flag.Parse()

	var benchmarkLabelMap = map[string]bench.Benchmark{
		"read":  &bench.ReadBenchmark{},
		"write": &bench.WriteBenchmark{},
		"sync":  &bench.SyncBenchmark{},
	}

	benchmark := benchmarkLabelMap[*benchmarkLabel]
	benchmarkResult := RunBenchmark(benchmark, *iterations, WithParallelism(*parallelism))
	fmt.Printf("%s benchmark = %+v\n", *benchmarkLabel, benchmarkResult)
}
