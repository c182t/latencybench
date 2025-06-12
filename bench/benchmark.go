package bench

import (
	"time"
)

type BenchmarkResult struct {
	Min time.Duration
	Avg time.Duration
	Max time.Duration
	P95 time.Duration
}
