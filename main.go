package main

import (
	"fmt"
	"latencybench/bench"
	"sort"
	"time"
)

func RunBenchmark(fn func() (time.Duration, error), iterations int) (bench.BenchmarkResult, error) {
	durations := make([]time.Duration, iterations)
	benchmarkResult := bench.BenchmarkResult{}
	var avg float64

	for i := 0; i < iterations; i++ {
		var err error
		durations[i], err = fn()

		if err != nil {
			return benchmarkResult, fmt.Errorf("RunBenchmark() - Failed to run benchmark at [%d] iteration; error: %v", i, err)
		}

		if i == 0 {
			benchmarkResult.Min = durations[i]
			benchmarkResult.Avg = durations[i]
			benchmarkResult.Max = durations[i]
		} else {
			if durations[i] < benchmarkResult.Min {
				benchmarkResult.Min = durations[i]
			}
			if durations[i] > benchmarkResult.Max {
				benchmarkResult.Max = durations[i]
			}
			avg += float64(benchmarkResult.Avg) / float64(iterations)
		}
	}

	benchmarkResult.Avg = time.Duration(avg)
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})
	benchmarkResult.P95 = durations[int(0.95*float64(len(durations)))]

	return benchmarkResult, nil
}

func main() {
	fmt.Println("Syscall Latency Benchmarks: ")

	rb := bench.ReadBenchmark{}
	rb.Open()
	readBenchmarkFn := func() (time.Duration, error) {
		return rb.RunBenchmark()
	}
	readDuration, err := RunBenchmark(readBenchmarkFn, 10000)
	rb.Close()

	if err != nil {
		fmt.Printf("Error occurred in BenchmarkRead: %v", err)
	} else {
		fmt.Printf("read duration = %v\n", readDuration)
	}
}
