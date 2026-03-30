package processSubmission_test

import (
	"testing"
)

func TestProcessSubmission_Node(t *testing.T) {
	tests := []ExecutionTestCase{
		{
			Name: "Node_Accepted",
			SourceCode: `const fs = require('fs');
const input = fs.readFileSync('/dev/stdin', 'utf-8').trim().split(/\s+/);
if (input.length >= 2) {
    console.log(parseInt(input[0]) + parseInt(input[1]));
}`,
			ExpectedState: "Accepted",
		},
		{
			Name: "Node_Wrong_Answer",
			SourceCode: `console.log("0");`,
			ExpectedState: "Wrong Answer",
		},
		{
			Name: "Node_Runtime_Error_Syntax",
			SourceCode: `const x = ;`, // Syntax error evaluated at runtime load. Node process terminates with exit code 1.
			ExpectedState: "Runtime Error",
		},
		{
			Name: "Node_Runtime_Error_Exception",
			SourceCode: `throw new Error("Crash");`,
			ExpectedState: "Runtime Error",
		},
		{
			Name: "Node_Time_Limit_Exceeded",
			SourceCode: `function delay(ms) { const start = Date.now(); while (Date.now() - start < ms) {} }
delay(5000);`,
			TimeLimitMs: 1500,
			ExpectedState: "Time Limit Exceeded",
		},
		{
			Name: "Node_Memory_Limit_Exceeded",
			SourceCode: `const arr = [];
while(true) {
    arr.push(new Array(1000).fill(1));
}`,
			MemoryLimitKb: 50 * 1024,
			ExpectedState: "Memory Limit Exceeded",
		},
	}

	tests = append(tests, GenerateVariations("Node_Bulk_Accepted", 
		`const fs = require('fs');
		const input = fs.readFileSync('/dev/stdin', 'utf-8').trim().split(/\s+/);
		console.log(parseInt(input[0]) + parseInt(input[1]));`, 
		"Accepted", 15)...)

	tests = append(tests, GenerateVariations("Node_Bulk_WA", 
		`console.log("-999");`, 
		"Wrong Answer", 10)...)

	runLanguageTests(t, "node", tests)
}
