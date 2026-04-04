package processSubmission

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rishik92/velox/judge"
	"github.com/rishik92/velox/runBatch"
)

type SystemError struct {
	msg string
}

func (e *SystemError) Error() string {
	return e.msg
}

type LanguageStrategy interface {
	Prepare(submissionID string, sourceCode string) (execCmd string, execArgs []string, filesToClean []string, err error)
}

var strategies = make(map[string]LanguageStrategy)

func RegisterStrategy(language string, strategy LanguageStrategy) {
	strategies[language] = strategy
}

type CSharpStrategy struct{}

func (s *CSharpStrategy) Prepare(submissionID string, sourceCode string) (string, []string, []string, error) {
	dirPath, dllPath, err := CompileInMemoryCSharp(submissionID, sourceCode)
	if err != nil {
		return "", nil, []string{dirPath}, err
	}
	return "dotnet", []string{dllPath}, []string{dirPath}, nil
}

type CStrategy struct{}

func (s *CStrategy) Prepare(submissionID string, sourceCode string) (string, []string, []string, error) {
	srcPath, binPath, err := CompileInMemoryC(submissionID, sourceCode)
	if err != nil {
		return "", nil, []string{srcPath, binPath}, err
	}
	return binPath, []string{}, []string{srcPath, binPath}, nil
}

type CPPStrategy struct{}

func (s *CPPStrategy) Prepare(submissionID string, sourceCode string) (string, []string, []string, error) {
	srcPath, binPath, err := CompileInMemoryCPP(submissionID, sourceCode)
	if err != nil {
		return "", nil, []string{srcPath, binPath}, err
	}
	return binPath, []string{}, []string{srcPath, binPath}, nil
}

type JavaStrategy struct{}

func (s *JavaStrategy) Prepare(submissionID string, sourceCode string) (string, []string, []string, error) {
	dirPath, className, err := CompileInMemoryJava(submissionID, sourceCode)
	if err != nil {
		return "", nil, []string{dirPath}, err
	}
	return "java", []string{"-cp", dirPath, className}, []string{dirPath}, nil
}

type PythonStrategy struct{}

func (s *PythonStrategy) Prepare(submissionID string, sourceCode string) (string, []string, []string, error) {
	scriptPath := filepath.Join(os.TempDir(), fmt.Sprintf("solution_%s.py", submissionID))
	if err := os.WriteFile(scriptPath, []byte(sourceCode), 0644); err != nil {
		return "", nil, []string{scriptPath}, &SystemError{msg: "System Error: Cannot write to RAM"}
	}
	return "python3", []string{scriptPath}, []string{scriptPath}, nil
}

type NodeStrategy struct{}

func (s *NodeStrategy) Prepare(submissionID string, sourceCode string) (string, []string, []string, error) {
	scriptPath := filepath.Join(os.TempDir(), fmt.Sprintf("solution_%s.js", submissionID))
	if err := os.WriteFile(scriptPath, []byte(sourceCode), 0644); err != nil {
		return "", nil, []string{scriptPath}, &SystemError{msg: "System Error: Cannot write to RAM"}
	}
	return "node", []string{scriptPath}, []string{scriptPath}, nil
}

type TSStrategy struct{}

func (s *TSStrategy) Prepare(submissionID string, sourceCode string) (string, []string, []string, error) {
	jsPath, tsPath, err := CompileInMemoryTS(submissionID, sourceCode)
	if err != nil {
		return "", nil, []string{tsPath, jsPath}, err
	}
	return "node", []string{jsPath}, []string{tsPath, jsPath}, nil
}

func init() {
	RegisterStrategy("csharp", &CSharpStrategy{})
	RegisterStrategy("c", &CStrategy{})
	RegisterStrategy("cpp", &CPPStrategy{})
	RegisterStrategy("java", &JavaStrategy{})
	RegisterStrategy("python", &PythonStrategy{})
	RegisterStrategy("node", &NodeStrategy{})
	RegisterStrategy("ts", &TSStrategy{})
}

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

	// 1. ROUTING & COMPILATION via Strategy Pattern
	strategy, exists := strategies[req.Language]
	if !exists {
		return judge.SubmissionResponse{OverallState: "Unsupported Language"}
	}

	cmd, args, cleanupFiles, err := strategy.Prepare(req.SubmissionID, req.SourceCode)
	filesToClean = append(filesToClean, cleanupFiles...)
	
	if err != nil {
		if sysErr, ok := err.(*SystemError); ok {
			return judge.SubmissionResponse{OverallState: sysErr.Error()}
		}
		return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
	}

	execCmd = cmd
	execArgs = args

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
    sourcePath := fmt.Sprintf("/dev/shm/solution_%s.ts", submissionID)
    jsPath := fmt.Sprintf("/dev/shm/solution_%s.js", submissionID)
    os.WriteFile(sourcePath, []byte(sourceCode), 0644)

    cmd := exec.Command("esbuild", sourcePath, "--outfile="+jsPath, "--platform=node", "--format=cjs")
    
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