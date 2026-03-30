package processSubmission_test

import (
	"testing"
)

func TestProcessSubmission_C(t *testing.T) {
	tests := []ExecutionTestCase{
		{
			Name: "C_Accepted",
			SourceCode: `#include <stdio.h>
int main() { int a, b; scanf("%d %d", &a, &b); printf("%d", a+b); return 0; }`,
			ExpectedState: "Accepted",
		},
		{
			Name: "C_Wrong_Answer",
			SourceCode: `#include <stdio.h>
int main() { int a, b; scanf("%d %d", &a, &b); printf("%d", a-b); return 0; }`, // Subtracts instead of adds
			ExpectedState: "Wrong Answer",
		},
		{
			Name: "C_Compile_Error",
			SourceCode: `#include <stdio.h>
int main() { int a, b scanf("%d %d", &a, &b); return 0; }`, // Missing semicolon
			ExpectedState: "Compile Error",
		},
		{
			Name: "C_Syntax_Error",
			SourceCode: `void fake_func() { printf("Hello"); `, // Missing closing brace
			ExpectedState: "Compile Error",
		},
		{
			Name: "C_Runtime_Error_Segfault",
			SourceCode: `#include <stdio.h>
int main() { int *p = NULL; *p = 10; return 0; }`,
			ExpectedState: "Runtime Error",
		},
		{
			Name: "C_Runtime_Error_Exit1",
			SourceCode: `#include <stdio.h>
#include <stdlib.h>
int main() { exit(1); return 0; }`,
			ExpectedState: "Runtime Error",
		},
		{
			Name: "C_Time_Limit_Exceeded",
			SourceCode: `#include <stdio.h>
int main() { while(1) {} return 0; }`,
			TimeLimitMs: 1000, 
			ExpectedState: "Time Limit Exceeded",
		},
		{
			Name: "C_Memory_Limit_Exceeded",
			SourceCode: `#include <stdio.h>
#include <stdlib.h>
#include <string.h>
int main() {
    int i;
    for (i = 0; i < 50000; i++) {
        void* ptr = malloc(1024 * 1024); // Allocate many MBs
		if (ptr == NULL) break;
		memset(ptr, 1, 1024 * 1024); // Touch memory physical allocation
    }
    return 0;
}`,
			MemoryLimitKb: 50 * 1024, // 50 MB
			ExpectedState: "Memory Limit Exceeded",
		},
	}

	// Generate bulk variations of the tests to meet rigorous load/count demands (total ~30 cases)
	tests = append(tests, GenerateVariations("C_Bulk_Accepted", 
		`#include <stdio.h>
		int main() { int a, b; scanf("%d %d", &a, &b); printf("%d", a+b); return 0; }`, 
		"Accepted", 15)...)

	tests = append(tests, GenerateVariations("C_Bulk_WA", 
		`#include <stdio.h>
		int main() { printf("-9999"); return 0; }`, 
		"Wrong Answer", 10)...)

	runLanguageTests(t, "c", tests)
}
