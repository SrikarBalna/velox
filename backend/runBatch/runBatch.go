package runBatch

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/rishik92/velox/judge"
)

func RunBatch(execCmd string, execArgs []string, testCases []judge.TestCase, timeLimitMs int, memoryLimitKb int) []judge.TestCaseResult {
	var results []judge.TestCaseResult

	for _, tc := range testCases {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeLimitMs)*time.Millisecond)
		
		// Setup the command (e.g., "python3 /dev/shm/solution_123.py" OR "/dev/shm/solution_123.out")
		cmd := exec.CommandContext(ctx, execCmd, execArgs...)

		// Pipe input from RAM directly to the process
		cmd.Stdin = strings.NewReader(tc.Input)

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		// --- RESOURCE TRACKING ---
		var timeMs int64
		var memoryKb int64
		if cmd.ProcessState != nil {
			timeMs = (cmd.ProcessState.UserTime() + cmd.ProcessState.SystemTime()).Milliseconds()
			if rusage, ok := cmd.ProcessState.SysUsage().(*syscall.Rusage); ok {
				// On Linux, Maxrss is reported in Kilobytes. On macOS, it is reported in Bytes!
				if runtime.GOOS == "darwin" {
					memoryKb = int64(rusage.Maxrss) / 1024
				} else {
					memoryKb = int64(rusage.Maxrss)
				}
			}
		}

		stderrStr := strings.TrimSpace(stderr.String())

		if memoryLimitKb > 0 && memoryKb > int64(memoryLimitKb) {
			results = append(results, judge.TestCaseResult{
				TestCaseID: tc.TestCaseID,
				Status: "Memory Limit Exceeded",
				Stderr:   stderrStr,
				TimeMs:   timeMs,
				MemoryKb: memoryKb,
			})
			cancel()
			continue
		}

		if ctx.Err() == context.DeadlineExceeded {
			results = append(results, judge.TestCaseResult{
				TestCaseID: tc.TestCaseID,
				Status: "Time Limit Exceeded",
				Stderr:   stderrStr,
				TimeMs:   timeMs,
				MemoryKb: memoryKb,
			})
			cancel()
			// break // Fail-Fast: Stop running the rest of the 20 cases
			continue
		}
		if err != nil {
			results = append(results, judge.TestCaseResult{
				TestCaseID: tc.TestCaseID, 
				Status: "Runtime Error",
				Stderr:   stderrStr,
				TimeMs:   timeMs,
				MemoryKb: memoryKb,
			})
			cancel()
			// break
			continue
		}

		actual := strings.TrimSpace(stdout.String())
		expected := strings.TrimSpace(tc.ExpectedOutput)

		if actual == expected {
			results = append(results, judge.TestCaseResult{
				TestCaseID: tc.TestCaseID, 
				Status: "Accepted",
				TimeMs: timeMs,
				MemoryKb: memoryKb,
			})
		} else {
			results = append(results, judge.TestCaseResult{
				TestCaseID: tc.TestCaseID, 
				Status: "Wrong Answer", 
				ActualOutput: actual, // Only send actual output back if it failed, helps with debugging!
				TimeMs: timeMs,
				MemoryKb: memoryKb,
			})
			cancel()
			// break
			continue
		}
		
		cancel()
	}

	return results
}