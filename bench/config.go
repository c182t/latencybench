package bench

import (
	"os"

	"gopkg.in/yaml.v3"
)

type BenchmarkSuite struct {
	Benchmarks []BenchmarkOptions `yaml:"benchmarks"`
}

type BenchmarkOptions struct {
	Benchmark   string `yaml:"benchmark"`
	Iterations  int    `yaml:"iterations"`
	Parallelism int    `yaml:"parallelism"`
	BlockSize   uint   `yaml:"block_size"`
	Stride      uint   `yaml:"stride"`
	RawSyscall  bool   `yaml:"raw_syscall"`
}

func LoadBenchmarkSuite(configPath string) (*BenchmarkSuite, error) {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var benchmarkSuite BenchmarkSuite
	err = yaml.Unmarshal(configData, &benchmarkSuite)
	if err != nil {
		return nil, err
	}

	return &benchmarkSuite, nil
}
