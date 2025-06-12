package main

import (
	"fmt"
	"latencybench/bench"
	"time"
)

func main() {
	fmt.Println("Syscall Latency Benchmarks: ")

	rb := bench.ReadBenchmark{}
	rb.Open()
	readBenchmarkFn := func() (time.Duration, error) {
		return rb.BenchmarkOnce()
	}
	readDuration, err := bench.RunBenchmark(readBenchmarkFn, 10000)
	rb.Close()

	if err != nil {
		fmt.Printf("Error occurred in RunBenchmark: %v", err)
	} else {
		fmt.Printf("read duration = %+v\n", readDuration)
	}
}
