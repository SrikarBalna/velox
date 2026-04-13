# 3. Sequence Diagrams

These diagrams show the **exact message flow** for every user-facing operation in Velox, tracing every function call, Redis operation, and HTTP exchange from start to finish.

---

## 3.1 Code Submission Flow (Happy Path — Compiled Language, e.g. C++)

This is the complete end-to-end flow from the user submitting code to receiving "Accepted".

```mermaid
sequenceDiagram
    autonumber
    participant Client as Client (Browser / curl)
    participant API as API Server<br/>cmd/api/main.go
    participant AuthMW as APIKey Auth
    participant Redis as Redis
    participant LogSvc as APILogService
    participant Worker as Worker<br/>cmd/worker/main.go
    participant PS as processSubmission<br/>.ProcessSubmission()
    participant Compiler as g++ (Compiler)
    participant RB as runBatch<br/>.RunBatch()
    participant Binary as Compiled Binary

    Note over Client,Binary: ── PHASE 1: SUBMISSION ──

    Client->>+API: POST /submit<br/>{ language: "cpp", source_code: "...", test_cases: [...] }
    API->>AuthMW: Validate Bearer Token
    AuthMW-->>API: Proceed (APIKeyID injected)
    API->>API: Validate request<br/>(TimeLimitMs ≤ 5000, MemoryLimitKb ≤ 512000)
    API->>API: Generate UUID<br/>req.SubmissionID = uuid.New()
    API->>Redis: LPUSH "submissions" <serialized JSON>
    Redis-->>API: OK
    API-->>LogSvc: Async Log(State: Pending)
    API-->>-Client: 202 Accepted<br/>{ "submission_id": "abc-123" }

    Note over Client,Binary: ── PHASE 2: PROCESSING ──

    Worker->>+Redis: BRPOP "submissions" 5s
    Redis-->>-Worker: <serialized JSON>
    Worker->>Worker: json.Unmarshal → SubmissionRequest
    Worker->>+PS: ProcessSubmission(req)

    PS->>PS: switch req.Language → case "cpp"
    PS->>PS: Write source to /tmp/solution_abc-123.cpp
    PS->>+Compiler: exec.Command("g++", src, "-O2", "-o", bin)
    Compiler-->>-PS: Compilation successful

    PS->>+RB: RunBatch(binPath, [], testCases, 3000, 256000)

    loop For each TestCase
        RB->>RB: context.WithTimeout(3000ms)
        RB->>+Binary: exec.CommandContext(binPath)<br/>stdin ← testCase.Input
        Binary-->>-RB: stdout, stderr, ProcessState
        RB->>RB: Measure time via ProcessState.UserTime() + SystemTime()
        RB->>RB: Measure memory via syscall.Rusage.Maxrss
        RB->>RB: Compare strings.TrimSpace(stdout) == strings.TrimSpace(expected)
        RB->>RB: Append TestCaseResult{Status: "Accepted"}
    end

    RB-->>-PS: []TestCaseResult
    PS->>PS: Aggregate: all "Accepted" → OverallState = "Accepted"
    PS->>PS: defer cleanup: os.RemoveAll(srcPath, binPath)
    PS-->>-Worker: SubmissionResponse

    Worker->>Worker: json.Marshal(response)
    Worker->>Redis: LPUSH "results:abc-123" <serialized JSON>

    Note over Client,Binary: ── PHASE 3: POLLING ──

    Client->>+API: GET /status?submission_id=abc-123
    API->>AuthMW: Validate Bearer Token
    AuthMW-->>API: Proceed (APIKeyID injected)
    API->>+Redis: BRPOP "results:abc-123" 1s
    Redis-->>-API: <serialized JSON>
    API-->>LogSvc: UpdateResult(State: Accepted)
    API-->>LogSvc: Async Log(/status request)
    API-->>-Client: 200 OK<br/>{ "submission_id": "abc-123", "overall_state": "Accepted", "results": [...] }
```

---

## 3.2 Submission Flow — Compile Error Path

When the user's code has a syntax error and fails compilation.

```mermaid
sequenceDiagram
    autonumber
    participant Client as Client
    participant API as API Server
    participant Redis as Redis
    participant Worker as Worker
    participant PS as processSubmission
    participant Compiler as Compiler (gcc/g++/javac)

    Client->>+API: POST /submit { language: "cpp", source_code: "bad code..." }
    API->>API: Generate UUID
    API->>Redis: LPUSH "submissions" <JSON>
    API-->>-Client: 202 Accepted { "submission_id": "xyz-456" }

    Worker->>+Redis: BRPOP "submissions" 5s
    Redis-->>-Worker: <JSON>
    Worker->>+PS: ProcessSubmission(req)

    PS->>PS: Write source to /tmp/solution_xyz-456.cpp
    PS->>+Compiler: exec.Command("g++", src, "-O2", "-o", bin)
    Compiler-->>-PS: Exit code 1 + stderr output

    PS->>PS: Return SubmissionResponse: OverallState = "Compile Error"<br/>CompileError = "..."
    PS->>PS: defer cleanup: os.RemoveAll(srcPath, binPath)
    PS-->>-Worker: SubmissionResponse

    Worker->>Redis: LPUSH "results:xyz-456" <JSON>

    Client->>+API: GET /status?submission_id=xyz-456
    API->>+Redis: BRPOP "results:xyz-456" 1s
    Redis-->>-API: <JSON>
    API-->>-Client: { "overall_state": "Compile Error", "compile_error": "error: syntax error..." }
```

---

## 3.3 Submission Flow — Pending (Still Processing)

When the client polls before the worker has finished.

```mermaid
sequenceDiagram
    autonumber
    participant Client as Client
    participant API as API Server
    participant Redis as Redis

    Client->>+API: GET /status?submission_id=abc-123
    API->>+Redis: BRPOP "results:abc-123" 1s
    Note right of Redis: 1 second timeout expires<br/>No result yet
    Redis-->>-API: nil (timeout)
    API-->>-Client: 200 OK<br/>{ "status": "pending" }

    Note over Client: Client waits 500ms and retries...

    Client->>+API: GET /status?submission_id=abc-123
    API->>+Redis: BRPOP "results:abc-123" 1s
    Redis-->>-API: <result JSON>
    API-->>-Client: 200 OK<br/>{ "overall_state": "Accepted", "results": [...] }
```

---

## 3.4 Submission Flow — Runtime Error / Time Limit Exceeded

```mermaid
sequenceDiagram
    autonumber
    participant RB as RunBatch
    participant Binary as User's Binary
    participant Ctx as context.WithTimeout

    RB->>Ctx: Create context with timeLimitMs deadline
    RB->>+Binary: exec.CommandContext(binary)<br/>stdin ← testCase.Input
    
    alt Normal exit with non-zero code
        Binary-->>-RB: exit code ≠ 0, stderr populated
        RB->>RB: Append TestCaseResult{Status: "Runtime Error", Stderr: "..."}
    
    else Context deadline exceeded
        Ctx->>Binary: Kill (SIGKILL)
        Binary-->>RB: Process killed
        RB->>RB: ctx.Err() == DeadlineExceeded
        RB->>RB: Append TestCaseResult{Status: "Time Limit Exceeded"}
    
    else Memory limit exceeded
        Binary-->>RB: Normal exit
        RB->>RB: Check rusage.Maxrss > memoryLimitKb
        RB->>RB: Append TestCaseResult{Status: "Memory Limit Exceeded"}
    
    else Wrong answer
        Binary-->>RB: Normal exit, stdout ≠ expected
        RB->>RB: Append TestCaseResult{Status: "Wrong Answer"}
    end
```

---

## 3.5 Interpreted Language Flow (Python / Node.js)

For interpreted languages, there is no compilation step — the source code is written directly to a temp file and executed.

```mermaid
sequenceDiagram
    autonumber
    participant PS as processSubmission
    participant FS as Filesystem (/tmp)
    participant RB as runBatch
    participant Runtime as python3 / node

    PS->>PS: switch req.Language → case "python"
    PS->>FS: os.WriteFile("/tmp/solution_xxx.py", sourceCode)
    FS-->>PS: OK

    PS->>+RB: RunBatch("python3", ["/tmp/solution_xxx.py"], testCases, timeLimit, memLimit)

    loop For each TestCase
        RB->>+Runtime: exec.CommandContext("python3", scriptPath)<br/>stdin ← input
        Runtime-->>-RB: stdout, stderr
        RB->>RB: Compare output, measure time & memory
    end

    RB-->>-PS: []TestCaseResult

    PS->>PS: defer: os.RemoveAll("/tmp/solution_xxx.py")
```

---

## 3.6 Explanation

### Key Design Patterns Visible in the Sequence Diagrams

| Pattern | Where | Why |
|---------|-------|-----|
| **Async Job Queue** | API → Redis → Worker | The API never blocks on code execution. It returns `202 Accepted` immediately, enabling the frontend to poll for results. This is essential because code execution can take seconds. |
| **Blocking Pop (BRPOP)** | Worker ← Redis | The worker uses `BRPOP` with a 5-second timeout to efficiently wait for new jobs without busy-polling. This is Redis's built-in mechanism for pub-sub-style work distribution. |
| **Context Timeout** | `RunBatch` | Each test case execution gets its own `context.WithTimeout` to enforce the time limit. If the child process exceeds it, Go's `exec.CommandContext` automatically sends `SIGKILL`. |
| **Deferred Cleanup** | `processSubmission.ProcessSubmission()` | All temp files (source code, compiled binaries) are added to a `filesToClean` slice. A `defer` block runs `os.RemoveAll()` on each, guaranteeing cleanup even if the function panics. |
| **Fail-through Evaluation** | `RunBatch` loop | The loop does NOT break on the first failure — it continues running all test cases so the user can see exactly which ones failed. A commented-out `break` shows this was a deliberate design decision. |
