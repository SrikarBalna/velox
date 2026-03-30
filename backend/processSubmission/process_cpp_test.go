package processSubmission_test

import (
	"testing"
)

func TestProcessSubmission_CPP(t *testing.T) {
	tests := []ExecutionTestCase{
		{
			Name: "CPP_Accepted",
			SourceCode: `#include <iostream>
using namespace std; 
int main() { int a, b; cin >> a >> b; cout << a+b; return 0; }`,
			ExpectedState: "Accepted",
		},
		{
			Name: "CPP_Wrong_Answer",
			SourceCode: `#include <iostream>
using namespace std; 
int main() { int a, b; cin >> a >> b; cout << a * b; return 0; }`, 
			ExpectedState: "Wrong Answer",
		},
		{
			Name: "CPP_Compile_Error",
			SourceCode: `#include <iostream>
int main() { cout << "Missing std namespace without using directive"; return 0; }`,
			ExpectedState: "Compile Error",
		},
		{
			Name: "CPP_Runtime_Error_Exception",
			SourceCode: `#include <iostream>
#include <vector>
using namespace std;
int main() { vector<int> v; cout << v.at(5); return 0; }`, // Out of bounds throws exception -> terminates abnormally
			ExpectedState: "Runtime Error",
		},
		{
			Name: "CPP_Time_Limit_Exceeded",
			SourceCode: `#include <iostream>
int main() { while(true) {} return 0; }`,
			TimeLimitMs: 1000, 
			ExpectedState: "Time Limit Exceeded",
		},
		{
			Name: "CPP_Memory_Limit_Exceeded",
			SourceCode: `#include <iostream>
#include <vector>
using namespace std;
int main() {
    vector<vector<int>> matrix;
    for (int i=0; i<10000; i++) {
        matrix.push_back(vector<int>(10000, 1));
    }
    return 0;
}`,
			MemoryLimitKb: 50 * 1024,
			ExpectedState: "Memory Limit Exceeded",
		},
	}

	tests = append(tests, GenerateVariations("CPP_Bulk_Accepted", 
		`#include <iostream>
		using namespace std; int main() { int a, b; if (cin >> a >> b) { cout << a+b; } return 0; }`, 
		"Accepted", 15)...)

	tests = append(tests, GenerateVariations("CPP_Bulk_WA", 
		`#include <iostream>
		using namespace std; int main() { cout << "-9999"; return 0; }`, 
		"Wrong Answer", 10)...)

	runLanguageTests(t, "cpp", tests)
}
