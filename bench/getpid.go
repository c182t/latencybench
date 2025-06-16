package bench

import (
	"os"
	"syscall"
	"time"
)

type GetPidBenchmark struct {
	Options *BenchmarkOptions
}

func (gpb *GetPidBenchmark) Setup() error {
	return nil
}

func (gpb *GetPidBenchmark) RunOnce() (time.Duration, error) {
	var duration time.Duration
	if gpb.GetOptions().RawSyscall {
		startTime := time.Now()
		const SYS_GETPID = 39 // x86_64 Linux
		r1, _, _ := syscall.RawSyscall(uintptr(SYS_GETPID), 0, 0, 0)
		duration = time.Since(startTime)
		r1 = r1 + 0
	} else {
		startTime := time.Now()
		pid := os.Getegid()
		duration = time.Since(startTime)
		pid = pid + 0
	}

	return duration, nil
}

func (gpb *GetPidBenchmark) Teardown() {}

func (gpb *GetPidBenchmark) Clone() Benchmark {
	clone := *gpb
	return &clone
}

func (gpb *GetPidBenchmark) GetOptions() *BenchmarkOptions {
	return gpb.Options
}
