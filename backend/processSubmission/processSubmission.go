package processSubmission

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rishik92/velox/judge"
	"github.com/rishik92/velox/runBatch"
)

func ProcessSubmission(req judge.SubmissionRequest) judge.SubmissionResponse {
	var execCmd string
	var execArgs []string

	// The `defer` block runs immediately before the function exits, NO MATTER WHAT (even on errors or panics).
	var filesToClean []string
	defer func() {
		for _, file := range filesToClean {
			os.RemoveAll(file)
		}
	}()

	// 1. ROUTING & COMPILATION
	switch req.Language {
	case "csharp":
		dirPath, dllPath, err := CompileInMemoryCSharp(req.SubmissionID, req.SourceCode)
		filesToClean = append(filesToClean, dirPath) // Add to cleanup list IMMEDIATELY
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = "dotnet"
		execArgs = []string{dllPath}

	case "c":
		srcPath, binPath, err := CompileInMemoryC(req.SubmissionID, req.SourceCode)
		filesToClean = append(filesToClean, srcPath, binPath) 
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = binPath
		execArgs = []string{}

	case "cpp":
		srcPath, binPath, err := CompileInMemoryCPP(req.SubmissionID, req.SourceCode)
		filesToClean = append(filesToClean, srcPath, binPath)
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = binPath
		execArgs = []string{}

	case "java":
		dirPath, className, err := CompileInMemoryJava(req.SubmissionID, req.SourceCode)
		filesToClean = append(filesToClean, dirPath)
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = "java"
		execArgs = []string{"-cp", dirPath, className}

	case "python":
		scriptPath := filepath.Join(os.TempDir(), fmt.Sprintf("solution_%s.py", req.SubmissionID))
		filesToClean = append(filesToClean, scriptPath)
		if err := os.WriteFile(scriptPath, []byte(req.SourceCode), 0644); err != nil {
			return judge.SubmissionResponse{OverallState: "System Error: Cannot write to RAM"}
		}
		execCmd = "python3"
		execArgs = []string{scriptPath}

	case "node":
		scriptPath := filepath.Join(os.TempDir(), fmt.Sprintf("solution_%s.js", req.SubmissionID))
		filesToClean = append(filesToClean, scriptPath)
		if err := os.WriteFile(scriptPath, []byte(req.SourceCode), 0644); err != nil {
			return judge.SubmissionResponse{OverallState: "System Error: Cannot write to RAM"}
		}
		execCmd = "node"
		execArgs = []string{scriptPath}

	case "ts":
		jsPath, tsPath, err := CompileInMemoryTS(req.SubmissionID, req.SourceCode)
		filesToClean = append(filesToClean, tsPath, jsPath)
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = "node"
		execArgs = []string{jsPath}

	default:
		return judge.SubmissionResponse{OverallState: "Unsupported Language"}
	}
	
	timeLimit := req.TimeLimitMs
	if timeLimit <= 0 {
		timeLimit = 3000
	}
	memLimit := req.MemoryLimitKb
	if memLimit <= 0 {
		memLimit = 256000
	}

	results := runBatch.RunBatch(execCmd, execArgs, req.TestCases, timeLimit, memLimit)

	// 3. AGGREGATE RESULTS (Cleanup is now handled safely by the `defer` block above)
	overallState := "Accepted"
	for _, res := range results {
		if res.Status != "Accepted" {
			overallState = res.Status
			// Do NOT break here if you want to show the user which specific test cases failed. 
			// If you only care about the first failure (Fail-Fast), you can uncomment the break.
		}
	}

	return judge.SubmissionResponse{
		SubmissionID: req.SubmissionID,
		OverallState: overallState,
		Results:      results,
	}
}

func CompileInMemoryC(submissionID, sourceCode string) (string, string, error) {
	sourcePath := filepath.Join(os.TempDir(), fmt.Sprintf("solution_%s.c", submissionID))
	binaryPath := filepath.Join(os.TempDir(), fmt.Sprintf("solution_%s_c", submissionID))
	os.WriteFile(sourcePath, []byte(sourceCode), 0644)

	cmd := exec.Command("gcc", sourcePath, "-O2", "-o", binaryPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return sourcePath, binaryPath, fmt.Errorf("compile error: %s", string(out)) // Cleaned up error formatting
	}
	return sourcePath, binaryPath, nil
}

func CompileInMemoryCPP(submissionID, sourceCode string) (string, string, error) {
	sourcePath := filepath.Join(os.TempDir(), fmt.Sprintf("solution_%s.cpp", submissionID))
	binaryPath := filepath.Join(os.TempDir(), fmt.Sprintf("solution_%s_cpp", submissionID))
	os.WriteFile(sourcePath, []byte(sourceCode), 0644)

	cmd := exec.Command("g++", sourcePath, "-O2", "-o", binaryPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return sourcePath, binaryPath, fmt.Errorf("compile error: %s", string(out))
	}
	return sourcePath, binaryPath, nil
}

func CompileInMemoryJava(submissionID, sourceCode string) (string, string, error) {
	dirPath := filepath.Join(os.TempDir(), fmt.Sprintf("sol_%s", submissionID))
	os.MkdirAll(dirPath, 0755)
	sourcePath := fmt.Sprintf("%s/Main.java", dirPath)
	os.WriteFile(sourcePath, []byte(sourceCode), 0644)

	cmd := exec.Command("javac", sourcePath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return dirPath, "Main", fmt.Errorf("compile error: %s", string(out))
	}
	return dirPath, "Main", nil
}

func CompileInMemoryTS(submissionID, sourceCode string) (string, string, error) {
	sourcePath := filepath.Join(os.TempDir(), fmt.Sprintf("solution_%s.ts", submissionID))
	jsPath := filepath.Join(os.TempDir(), fmt.Sprintf("solution_%s.js", submissionID))
	os.WriteFile(sourcePath, []byte(sourceCode), 0644)

	// Use tsc for TS to JS compilation with type checking
	cmd := exec.Command("npx", "tsc", sourcePath, "--skipLibCheck")
	if out, err := cmd.CombinedOutput(); err != nil {
		return jsPath, sourcePath, fmt.Errorf("compile error: %s", string(out))
	}
	return jsPath, sourcePath, nil
}

func CompileInMemoryCSharp(submissionID, sourceCode string) (string, string, error) {
	dirPath := filepath.Join(os.TempDir(), fmt.Sprintf("sol_cs_%s", submissionID))
	os.MkdirAll(dirPath, 0755)

	csprojContent := `<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>net8.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
    <Nullable>enable</Nullable>
  </PropertyGroup>
</Project>`
	
	os.WriteFile(dirPath+"/project.csproj", []byte(csprojContent), 0644)
	os.WriteFile(dirPath+"/Program.cs", []byte(sourceCode), 0644)

	cmd := exec.Command("dotnet", "build", dirPath, "-c", "Release", "-o", dirPath+"/out", "-v", "q", "--nologo")
	if out, err := cmd.CombinedOutput(); err != nil {
		return dirPath, dirPath + "/out/project.dll", fmt.Errorf("compile error: %s", string(out))
	}

	return dirPath, dirPath + "/out/project.dll", nil
}