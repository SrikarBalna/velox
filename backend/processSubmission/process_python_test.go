package processSubmission_test

import (
	"testing"
)

func TestProcessSubmission_Python(t *testing.T) {
	tests := []ExecutionTestCase{
		{
			Name:          "Py_Accepted",
			SourceCode:    `import sys; [print(sum(map(int, line.split()))) for line in sys.stdin]`,
			ExpectedState: "Accepted",
		},
		{
			Name: "Py_Accepted_MultiLine",
			SourceCode: `a, b = map(int, input().split())
print(a + b)`,
			ExpectedState: "Accepted",
		},
		{
			Name:          "Py_Wrong_Answer",
			SourceCode:    `a, b = map(int, input().split()); print(a * b)`,
			ExpectedState: "Wrong Answer",
		},
		{
			Name: "Py_Runtime_Error_NameError",
			SourceCode: `a, b = map(int, input().split())
prin(a+b)`,  // Typo prin instead of print
			ExpectedState: "Runtime Error",
		},
		{
			Name: "Py_Runtime_Error_Div0",
			SourceCode: `print(1/0)`,
			ExpectedState: "Runtime Error",
		},
		{
			Name: "Py_Syntax_Error",
			SourceCode: `if True
    print("No colon")`,
			ExpectedState: "Runtime Error", // Python compilation errors output to stderr on execution, so it triggers RE or WA depending on exact routing.
		},
		{
			Name: "Py_Time_Limit_Exceeded",
			SourceCode: `import time; time.sleep(10)`,
			TimeLimitMs:   1000, 
			ExpectedState: "Time Limit Exceeded",
		},
		{
			Name: "Py_Memory_Limit_Exceeded",
			SourceCode: `arr = [0] * (10**8)`, // Massive array 
			MemoryLimitKb: 50 * 1024, // 50MB limits
			ExpectedState: "Memory Limit Exceeded",
		},
	}

	// Generate bulk volume validation
	tests = append(tests, GenerateVariations("Py_Bulk_Accepted", 
		`a, b = map(int, input().split()); print(a + b)`, 
		"Accepted", 15)...)

	tests = append(tests, GenerateVariations("Py_Bulk_WA", 
		`print("-9999")`, 
		"Wrong Answer", 10)...)

	runLanguageTests(t, "python", tests)
}
