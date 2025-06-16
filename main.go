package main

import (
	"flag"
	"fmt"
	"latencybench/bench"
	"runtime"

	"gopkg.in/yaml.v3"
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
	fmt.Printf("Starting latencybench with GOMAXPROCS=%d\n", runtime.GOMAXPROCS(0))

	var benchmarkLabelMap = map[string]func(options bench.BenchmarkOptions) bench.Benchmark{
		"read":  func(options bench.BenchmarkOptions) bench.Benchmark { return &bench.ReadBenchmark{Options: &options} },
		"write": func(options bench.BenchmarkOptions) bench.Benchmark { return &bench.WriteBenchmark{Options: &options} },
		"sync":  func(options bench.BenchmarkOptions) bench.Benchmark { return &bench.SyncBenchmark{Options: &options} },
		"memory_copy": func(options bench.BenchmarkOptions) bench.Benchmark {
			return &bench.MemoryCopyBenchmark{Options: &options}
		},
		"memory_stride": func(options bench.BenchmarkOptions) bench.Benchmark {
			return &bench.MemoryStrideBenchmark{Options: &options}
		},
		"getpid": func(options bench.BenchmarkOptions) bench.Benchmark {
			return &bench.GetPidBenchmark{Options: &options}
		},
		"retint": func(options bench.BenchmarkOptions) bench.Benchmark {
			return &bench.RetIntBenchmark{Options: &options}
		},
	}

	configPath := flag.String("config_path", "", "yaml config file path")
	benchmarkLabel := flag.String("benchmark", "read", "Benchmark to execute (e.g. read, write, sync)")
	iterations := flag.Int("iterations", 1000, "Number of benchmark iterations")
	parallelism := flag.Int("parallelism", 1, "Number of threads to run in parallel")
	blockSize := flag.Uint("block_size", 4*1024, "Block size (e.g. 4096)")
	stride := flag.Uint("stride", 64, "Memory access stride (e.g. 1 for sequential, 64 to jump 64 bytes, etc)")
	rawSyscall := flag.Bool("raw_syscall", false, "Use RawSyscal (Linux x86_64 only)")

	flag.Parse()

	if *configPath != "" {
		suite, err := bench.LoadBenchmarkSuite(*configPath)
		if err != nil {
			fmt.Printf("Failed to load config [%s]: %v", *configPath, err)
		}

		var resultsOptions []interface{}
		for _, options := range suite.Benchmarks {
			benchmark := benchmarkLabelMap[options.Benchmark]
			if benchmark == nil {
				fmt.Printf("Unknown benchmark: %s", options.Benchmark)
				continue
			}

			result := RunBenchmark(benchmark(options), options.Iterations)

			combinedResult := map[string]interface{}{
				"options":           options,
				"aggregated_result": result,
			}

			resultsOptions = append(resultsOptions, combinedResult)
		}

		yamlResultsOptions, _ := yaml.Marshal(resultsOptions)
		fmt.Println(string(yamlResultsOptions))

	} else {
		options := bench.BenchmarkOptions{Benchmark: *benchmarkLabel,
			Iterations:  *iterations,
			Parallelism: *parallelism,
			BlockSize:   *blockSize,
			Stride:      *stride,
			RawSyscall:  *rawSyscall}

		benchmark := benchmarkLabelMap[*benchmarkLabel]
		benchmarkResult := RunBenchmark(benchmark(options),
			*iterations)

		fmt.Printf("%s benchmark = %+v\n", *benchmarkLabel, benchmarkResult)
	}
}
