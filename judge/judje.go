package judge

// Incoming JSON from the Frontend API
type SubmissionRequest struct {
	SubmissionID string     `json:"submission_id"`
	Language     string     `json:"language"` // e.g., "cpp", "python"
	SourceCode   string     `json:"source_code"`
	TestCases    []TestCase `json:"test_cases"` // Up to 20 cases
}

type TestCase struct {
	TestCaseID     int    `json:"test_case_id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

// Outgoing JSON back to the Frontend Dashboard
type SubmissionResponse struct {
	SubmissionID string           `json:"submission_id"`
	OverallState string           `json:"overall_state"` // "Accepted", "Wrong Answer", "Compile Error"
	CompileError string           `json:"compile_error,omitempty"`
	Results      []TestCaseResult `json:"results"`
}

type TestCaseResult struct {
	TestCaseID   int    `json:"test_case_id"`
	Status       string `json:"status"` // "Passed", "Wrong Answer", "Runtime Error"
	ActualOutput string `json:"actual_output,omitempty"`
	Stderr       string `json:"stderr,omitempty"`
	TimeMs       int64  `json:"time_ms,omitempty"`
	MemoryKb     int64  `json:"memory_kb,omitempty"`
}