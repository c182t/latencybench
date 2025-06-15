package main

import (
	"flag"
	"fmt"
	"latencybench/bench"
)

func RunBenchmark(b bench.Benchmark, n int) bench.BenchmarkAggregatedResult {
	var err error
	var benchmarkAggResult bench.BenchmarkAggregatedResult

	switch {
	case b.GetOptions().Parallelism == 1:
		benchmarkAggResult, err = bench.RunBenchmarkSerial(b)
	case b.GetOptions().Parallelism > 1:
		benchmarkAggResult, err = bench.RunBenchmarkParallel(b)
	default:
		fmt.Printf("Error occurred in RunBenchmark - incorrect parallelism value: %d\n", b.GetOptions().Parallelism)
		err = fmt.Errorf("invalid parallelism value: %d", b.GetOptions().Parallelism)
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

	options := bench.BenchmarkOptions{Benchmark: *benchmarkLabel,
		Iterations:  *iterations,
		Parallelism: *parallelism,
		BlockSize:   *blockSize}

	var benchmarkLabelMap = map[string]bench.Benchmark{
		"read":        &bench.ReadBenchmark{Options: &options},
		"write":       &bench.WriteBenchmark{Options: &options},
		"sync":        &bench.SyncBenchmark{Options: &options},
		"memory_copy": &bench.MemoryCopyBenchmark{Options: &options},
	}

	benchmark := benchmarkLabelMap[*benchmarkLabel]
	benchmarkResult := RunBenchmark(benchmark,
		*iterations)

	fmt.Printf("%s benchmark = %+v\n", *benchmarkLabel, benchmarkResult)
}
