package processSubmission_test

import (
	"testing"
)

func TestProcessSubmission_TS(t *testing.T) {
	tests := []ExecutionTestCase{
		{
			Name: "TS_Accepted",
			SourceCode: `// @ts-ignore
import * as fs from 'fs';
const input: string[] = fs.readFileSync('/dev/stdin', 'utf-8').trim().split(/\s+/);
if (input.length >= 2) {
    const a: number = parseInt(input[0]);
    const b: number = parseInt(input[1]);
    console.log(a + b);
}`,
			ExpectedState: "Accepted",
		},
		{
			Name: "TS_Compile_Error",
			SourceCode: `const a: number = "string";`, // TypeScript type error: Type 'string' is not assignable to type 'number'
			ExpectedState: "Compile Error",
		},
		{
			Name: "TS_Wrong_Answer",
			SourceCode: `console.log("0");`,
			ExpectedState: "Wrong Answer",
		},
		{
			Name: "TS_Runtime_Error_Exception",
			SourceCode: `throw new Error("Crash");`,
			ExpectedState: "Runtime Error",
		},
		{
			Name: "TS_Time_Limit_Exceeded",
			SourceCode: `const start = Date.now(); while (Date.now() - start < 5000) {}`,
			TimeLimitMs: 1500,
			ExpectedState: "Time Limit Exceeded",
		},
	}

	tests = append(tests, GenerateVariations("TS_Bulk_Accepted", 
		`// @ts-ignore
import * as fs from 'fs';
		const input: string[] = fs.readFileSync('/dev/stdin', 'utf-8').trim().split(/\s+/);
		console.log(parseInt(input[0]) + parseInt(input[1]));`, 
		"Accepted", 15)...)

	runLanguageTests(t, "ts", tests)
}
