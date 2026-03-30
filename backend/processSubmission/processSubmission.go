package processSubmission

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rishik92/velox/judge"
	"github.com/rishik92/velox/runBatch"
)

func ProcessSubmission(req judge.SubmissionRequest) judge.SubmissionResponse {
	var execCmd string
	var execArgs []string
	
	var filesToClean []string

	// 1. ROUTING: Compiled vs Interpreted
	switch req.Language {
	case "c":
		binaryPath, err := CompileInMemoryC(req.SubmissionID, req.SourceCode)
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = binaryPath
		execArgs = []string{}
		filesToClean = append(filesToClean, binaryPath)

	case "cpp":
		binaryPath, err := CompileInMemoryCPP(req.SubmissionID, req.SourceCode)
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = binaryPath
		execArgs = []string{}
		filesToClean = append(filesToClean, binaryPath)

	case "java":
		// Java requires defining a class name, usually Main. We create a directory for this submission.
		dirPath, className, err := CompileInMemoryJava(req.SubmissionID, req.SourceCode)
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = "java"
		execArgs = []string{"-cp", dirPath, className}
		filesToClean = append(filesToClean, dirPath) // Clean up the entire submission directory

	case "python":
		scriptPath := fmt.Sprintf("%s/solution_%s.py", getTempDir(), req.SubmissionID)
		if err := os.WriteFile(scriptPath, []byte(req.SourceCode), 0644); err != nil {
			return judge.SubmissionResponse{OverallState: "System Error: Cannot write to RAM"}
		}
		execCmd = "python3"
		execArgs = []string{scriptPath}
		filesToClean = append(filesToClean, scriptPath)

	case "node":
		scriptPath := fmt.Sprintf("%s/solution_%s.js", getTempDir(), req.SubmissionID)
		if err := os.WriteFile(scriptPath, []byte(req.SourceCode), 0644); err != nil {
			return judge.SubmissionResponse{OverallState: "System Error: Cannot write to RAM"}
		}
		execCmd = "node"
		execArgs = []string{scriptPath}
		filesToClean = append(filesToClean, scriptPath)

	case "ts":
		jsPath, tsPath, err := CompileInMemoryTS(req.SubmissionID, req.SourceCode)
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = "node"
		execArgs = []string{jsPath}
		filesToClean = append(filesToClean, tsPath, jsPath)

	default:
		return judge.SubmissionResponse{OverallState: "Unsupported Language"}
	}

	// 2. EXECUTION: Run the batch with the prepared command
	timeLimit := req.TimeLimitMs
	if timeLimit <= 0 {
		timeLimit = 3000 // 3 seconds default
	}
	memLimit := req.MemoryLimitKb
	if memLimit <= 0 {
		memLimit = 256000 // 256MB default
	}

	results := runBatch.RunBatch(execCmd, execArgs, req.TestCases, timeLimit, memLimit)

	// 3. CLEANUP: Delete files from RAM-disk or Temp disk
	defer func() {
		for _, file := range filesToClean {
			_ = os.RemoveAll(file)
		}
	}()

	// 4. AGGREGATE RESULTS
	overallState := "Accepted"
	for _, res := range results {
		if res.Status != "Accepted" {
			overallState = res.Status // e.g., "Wrong Answer" or "Time Limit Exceeded"
			// break // Fail-Fast: We found an error, no need to check the rest
		}
	}

	return judge.SubmissionResponse{
		SubmissionID: req.SubmissionID,
		OverallState: overallState,
		Results:      results,
	}
}

func CompileInMemoryC(submissionID, sourceCode string) (string, error) {
	sourcePath := fmt.Sprintf("%s/solution_%s.c", getTempDir(), submissionID)
	binaryPath := fmt.Sprintf("%s/solution_%s", getTempDir(), submissionID)
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0644); err != nil {
		return "", fmt.Errorf("failed to write source: %v", err)
	}

	cmd := exec.Command("gcc", sourcePath, "-o", binaryPath, "-O2", "-Wall", "-lm")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("compile error: %v, %s", err, string(out))
	}
	return binaryPath, nil
}

func CompileInMemoryCPP(submissionID, sourceCode string) (string, error) {
	sourcePath := fmt.Sprintf("%s/solution_%s.cpp", getTempDir(), submissionID)
	binaryPath := fmt.Sprintf("%s/solution_%s_cpp", getTempDir(), submissionID)
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0644); err != nil {
		return "", fmt.Errorf("failed to write source: %v", err)
	}

	cmd := exec.Command("g++", sourcePath, "-o", binaryPath, "-O2", "-Wall")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("compile error: %v, %s", err, string(out))
	}
	return binaryPath, nil
}

func CompileInMemoryJava(submissionID, sourceCode string) (string, string, error) {
	dirPath := fmt.Sprintf("%s/sol_%s", getTempDir(), submissionID)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create dir: %v", err)
	}
	sourcePath := fmt.Sprintf("%s/Main.java", dirPath)
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write source: %v", err)
	}

	cmd := exec.Command("javac", sourcePath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", "", fmt.Errorf("compile error: %v, %s", err, string(out))
	}
	return dirPath, "Main", nil
}

func CompileInMemoryTS(submissionID, sourceCode string) (string, string, error) {
	sourcePath := fmt.Sprintf("%s/solution_%s.ts", getTempDir(), submissionID)
	jsPath := fmt.Sprintf("%s/solution_%s.js", getTempDir(), submissionID)
	if err := os.WriteFile(sourcePath, []byte(sourceCode), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write source: %v", err)
	}

	// --skipLibCheck prevents compilation from crashing due to random @types installed globally or in parent folders
	cmd := exec.Command("npx", "tsc", sourcePath, "--target", "ES2022", "--module", "commonjs", "--esModuleInterop", "--skipLibCheck")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", "", fmt.Errorf("compile error: exit status %v, %s", err, string(out))
	}
	return jsPath, sourcePath, nil
}

func getTempDir() string {
	if _, err := os.Stat("/dev/shm"); err == nil {
		return "/dev/shm"
	}
	return os.TempDir()
}