package bench

type BenchmarkConfig struct {
	Benchmark   string
	Iterations  int
	Parallelism int
	BlockSize   uint
}

type BenchmarkOption func(*BenchmarkConfig)

func WithBlockSize(blockSize uint) BenchmarkOption {
	return func(bc *BenchmarkConfig) {
		bc.BlockSize = blockSize
	}
}

func WithBenchmark(benchmark string) BenchmarkOption {
	return func(bc *BenchmarkConfig) {
		bc.Benchmark = benchmark
	}
}

func WithIterations(iterations int) BenchmarkOption {
	return func(bc *BenchmarkConfig) {
		bc.Iterations = iterations
	}
}

func WithParallelism(parallelism int) BenchmarkOption {
	return func(bc *BenchmarkConfig) {
		bc.Parallelism = parallelism
	}
}
