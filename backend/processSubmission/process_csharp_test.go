package processSubmission_test

import (
	"testing"
)

func TestProcessSubmission_CSharp(t *testing.T) {
	tests := []ExecutionTestCase{
		{
			Name: "CS_Accepted",
			SourceCode: `using System;
class Program {
    static void Main() {
        string line;
        while ((line = Console.ReadLine()) != null) {
            var parts = line.Split(' ', StringSplitOptions.RemoveEmptyEntries);
            if (parts.Length == 2) {
                int a = int.Parse(parts[0]);
                int b = int.Parse(parts[1]);
                Console.WriteLine(a + b);
            }
        }
    }
}`,
			ExpectedState: "Accepted",
		},
		{
			Name: "CS_Wrong_Answer",
			SourceCode: `using System;
class Program {
    static void Main() {
        string line = Console.ReadLine();
        Console.WriteLine("-999");
    }
}`,
			ExpectedState: "Wrong Answer",
		},
		{
			Name: "CS_Runtime_Error_Exception",
			SourceCode: `using System;
class Program {
    static void Main() {
        throw new Exception("Test Exception");
    }
}`,
			ExpectedState: "Runtime Error",
		},
		{
			Name: "CS_Time_Limit_Exceeded",
			SourceCode: `using System;
class Program {
    static void Main() {
        System.Threading.Thread.Sleep(10000);
    }
}`,
			TimeLimitMs:   1000,
			ExpectedState: "Time Limit Exceeded",
		},
	}

	// Generate bulk variations
	tests = append(tests, GenerateVariations("CS_Bulk_Accepted",
		`using System;
class Program {
    static void Main() {
        string line = Console.ReadLine();
        if (line != null) {
            var parts = line.Split(' ');
            Console.WriteLine(int.Parse(parts[0]) + int.Parse(parts[1]));
        }
    }
}`,
		"Accepted", 10)...)

	runLanguageTests(t, "csharp", tests)
}
