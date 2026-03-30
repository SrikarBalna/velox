package processSubmission_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/rishik92/velox/judge"
	"github.com/rishik92/velox/processSubmission"
)

// ExecutionTestCase encapsulates all requirements for simulating an E2E submission.
type ExecutionTestCase struct {
	Name          string
	SourceCode    string
	TestCases     []judge.TestCase // If nil, a default 1+2=3 test case is injected
	TimeLimitMs   int
	MemoryLimitKb int
	ExpectedState string // e.g., "Accepted", "Compile Error", "Time Limit Exceeded"
}

// GenerateVariations is a helper that duplicates a base core logic across identical repetitive tests
// strictly to validate consistency and volume under different inputs as requested.
func GenerateVariations(baseName string, sourceCode string, expectedState string, count int) []ExecutionTestCase {
	var results []ExecutionTestCase
	for i := 1; i <= count; i++ {
		results = append(results, ExecutionTestCase{
			Name:          fmt.Sprintf("%s_variant_%d", baseName, i),
			SourceCode:    sourceCode,
			ExpectedState: expectedState,
			TestCases: []judge.TestCase{
				{TestCaseID: i, Input: fmt.Sprintf("%d %d\n", i, i*2), ExpectedOutput: fmt.Sprintf("%d", i+(i*2))},
			},
		})
	}
	return results
}

// runLanguageTests serves as the Table-Driven-Test runner for any given language suite.
func runLanguageTests(t *testing.T, language string, tests []ExecutionTestCase) {
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%s/%d_%s", language, i, tc.Name), func(t *testing.T) {
			
			// Always unique ID to ensure isolated OS file structures
			req := judge.SubmissionRequest{
				SubmissionID:  fmt.Sprintf("tsub_%s_%d_%d", language, time.Now().UnixNano(), rand.Intn(10000)),
				Language:      language,
				SourceCode:    tc.SourceCode,
				TestCases:     tc.TestCases,
				TimeLimitMs:   tc.TimeLimitMs,
				MemoryLimitKb: tc.MemoryLimitKb,
			}

			// Provide defaults
			if req.TimeLimitMs == 0 {
				req.TimeLimitMs = 2000 // 2s
			}
			if req.MemoryLimitKb == 0 {
				req.MemoryLimitKb = 256000 // 256MB
			}
			if len(req.TestCases) == 0 {
				req.TestCases = []judge.TestCase{
					{TestCaseID: 1, Input: "1 2\n", ExpectedOutput: "3"},
					{TestCaseID: 2, Input: "-5 10\n", ExpectedOutput: "5"},
				}
			}

			resp := processSubmission.ProcessSubmission(req)

			if resp.OverallState != tc.ExpectedState {
				t.Errorf("\nFAIL: Expected OverallState %q, got %q\nCompiler Error: %v\nResults Dump: %v", 
					tc.ExpectedState, resp.OverallState, resp.CompileError, resp.Results)
			}
		})
	}
}
