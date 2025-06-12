package main

import (
	"fmt"
	"latencybench/bench"
)

func main() {
	fmt.Println("Syscall Latency Benchmarks: ")

	readDuration, err := bench.BenchmarkRead()

	if err != nil {
		fmt.Printf("Error occurred in BenchmarkRead: %v", err)
	} else {
		fmt.Printf("read duration = %v\n", readDuration)
	}
}
