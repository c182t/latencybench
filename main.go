package main

import (
	"fmt"
	"latencybench/bench"
	"time"
)

/*
func RunReadBenchmark() bench.BenchmarkResult {
	rb := bench.ReadBenchmark{}
	rb.Open()
	defer rb.Close()

	readBenchmarkFn := func() (time.Duration, error) {
		return rb.BenchmarkOnce()
	}
	benchmarkResult, err := bench.RunBenchmark(readBenchmarkFn, 10000)

	if err != nil {
		fmt.Printf("Error occurred in RunBenchmark: %v", err)
	}

	return benchmarkResult
}
*/

func RunBenchmark(b bench.Benchmark, n int) bench.BenchmarkResult {
	b.Open()
	defer b.Close()

	readBenchmarkFn := func() (time.Duration, error) {
		return b.BenchmarkOnce()
	}
	benchmarkResult, err := bench.RunBenchmark(readBenchmarkFn, n)

	if err != nil {
		fmt.Printf("Error occurred in RunBenchmark: %v", err)
	}

	return benchmarkResult
}

func main() {
	fmt.Println("Syscall Latency Benchmarks: ")

	readBenchmarkResult := RunBenchmark(&bench.ReadBenchmark{}, 100000)
	writeBenchmarkResult := RunBenchmark(&bench.WriteBenchmark{}, 100000)

	fmt.Printf("read benchmark = %+v\n", readBenchmarkResult)
	fmt.Printf("write benchmark = %+v\n", writeBenchmarkResult)

}
