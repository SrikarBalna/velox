package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rishik92/velox/judge"
	"github.com/rishik92/velox/processSubmission"
	veloxRedis "github.com/rishik92/velox/shared/redis"
)

func main() {
	// 1. Connect to Redis
	veloxRedis.Connect()
	fmt.Println("Velox worker started. Waiting for submissions on queue 'submissions'...")

	// 2. Continuously poll Redis for submissions
	service := processSubmission.NewSubmissionService(&processSubmission.DefaultRunner{}, processSubmission.NewDefaultRegistry(), &processSubmission.LocalTempStorageAdapter{})

	for {
		raw, found := veloxRedis.PopSubmission("submissions", 5*time.Second)
		if !found {
			continue
		}

		// 3. Parse the JSON submission
		var req judge.SubmissionRequest
		if err := json.Unmarshal([]byte(raw), &req); err != nil {
			log.Printf("Invalid submission JSON: %v", err)
			continue
		}

		fmt.Printf("\n==================================\n")
		fmt.Printf("Processing %s (%s)...\n", req.SubmissionID, req.Language)
		fmt.Printf("==================================\n")

		// 4. Process the submission
		response := service.ProcessSubmission(req)

		// 5. Marshal and push result back to Redis
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("JSON marshaling failed: %v", err)
			continue
		}

		resultQueue := "results:" + req.SubmissionID
		if err := veloxRedis.PushResult(resultQueue, string(responseJSON)); err != nil {
			log.Printf("Failed to push result to Redis: %v", err)
		}

		// Also print to stdout for logging
		pretty, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(pretty))
	}
}

/*
// --- MOCK TESTS (commented out) ---
func runMockTests() {
	testCases := []judge.TestCase{
		{TestCaseID: 1, Input: "1 2\n", ExpectedOutput: "3"},
		{TestCaseID: 2, Input: "10 20\n", ExpectedOutput: "30"},
		{TestCaseID: 3, Input: "-5 5\n", ExpectedOutput: "0"},
	}

	mocks := []judge.SubmissionRequest{
		{SubmissionID: "mock_c", Language: "c", SourceCode: `#include <stdio.h>
int main() { int a,b; while(scanf("%d %d",&a,&b)!=EOF) printf("%d\n",a+b); return 0; }`, TestCases: testCases},
		{SubmissionID: "mock_cpp", Language: "cpp", SourceCode: `#include <iostream>
using namespace std; int main() { int a,b; while(cin>>a>>b) cout<<a+b<<endl; return 0; }`, TestCases: testCases},
		{SubmissionID: "mock_python", Language: "python", SourceCode: `import sys
for line in sys.stdin:
    if line.strip():
        a, b = map(int, line.split())
        print(a + b)
`, TestCases: testCases},
	}

	service := processSubmission.NewSubmissionService(&processSubmission.DefaultRunner{}, processSubmission.NewDefaultRegistry(), &processSubmission.LocalTempStorageAdapter{})
	for _, req := range mocks {
		fmt.Printf("\nProcessing %s (%s)...\n", req.SubmissionID, req.Language)
		response := service.ProcessSubmission(req)
		responseJSON, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(responseJSON))
	}
}
*/
