# 2. Relationship Diagram

This document shows **how every package and struct depends on the others** — both at the package level (import graph) and at the struct level (data flow).


## 2.1 Struct-Level Relationship Diagram

This shows how data structures relate to each other and flow through the system.

```mermaid
classDiagram
    direction LR

    class SubmissionRequest {
        +string SubmissionID
        +string Language
        +string SourceCode
        +TestCase[] TestCases
        +int TimeLimitMs
        +int MemoryLimitKb
    }

    class TestCase {
        +int TestCaseID
        +string Input
        +string ExpectedOutput
    }

    class SubmissionResponse {
        +string SubmissionID
        +string OverallState
        +string CompileError
        +TestCaseResult[] Results
    }

    class TestCaseResult {
        +int TestCaseID
        +string Status
        +string ActualOutput
        +string Input
        +string ExpectedOutput
        +string Stderr
        +int64 TimeMs
        +int64 MemoryKb
    }

    class ProcessSubmissionPkg {
        <<package: processSubmission>>
    }

    class RunBatchPkg {
        <<package: runBatch>>
    }

    class RedisClient {
        <<package: shared/redis>>
    }

    class APIServer {
        <<package: cmd/api>>
    }

    class WorkerProcess {
        <<package: cmd/worker>>
    }

    SubmissionRequest "1" *-- "1..*" TestCase : contains

    SubmissionResponse "1" *-- "0..*" TestCaseResult : contains

    APIServer ..> SubmissionRequest : deserializes from JSON
    APIServer ..> RedisClient : pushes SubmissionRequest JSON

    WorkerProcess ..> RedisClient : pops SubmissionRequest JSON
    WorkerProcess ..> SubmissionRequest : deserializes
    WorkerProcess ..> ProcessSubmissionPkg : calls ProcessSubmission()
    WorkerProcess ..> SubmissionResponse : serializes to JSON
    WorkerProcess ..> RedisClient : pushes SubmissionResponse JSON

    ProcessSubmissionPkg ..> SubmissionRequest : reads
    ProcessSubmissionPkg ..> RunBatchPkg : calls RunBatch()
    ProcessSubmissionPkg ..> SubmissionResponse : returns

    RunBatchPkg ..> TestCase : iterates over
    RunBatchPkg ..> TestCaseResult : produces

    APIServer ..> RedisClient : pops SubmissionResponse JSON
```

---

## 2.2 Frontend Component Hierarchy

```mermaid
graph TD
    subgraph "Next.js Frontend"
        Layout["RootLayout<br/><i>app/layout.tsx</i>"]
        HomePage["Home Page<br/><i>app/page.tsx</i>"]
        DocsPage["Docs Page<br/><i>app/docs/page.tsx</i>"]
        LoginPage["Login Page<br/><i>app/login/page.tsx</i>"]
        SignupPage["Signup Page<br/><i>app/signup/page.tsx</i>"]

        Navbar["Navbar"]
        Hero["Hero"]
        Comparison["Comparison"]
        Features["Features"]
        Footer["Footer"]
        SearchModal["SearchModal"]
        Button["Button"]
        Sidebar["Sidebar"]
        CodeBlock["CodeBlock"]
    end

    Layout --> HomePage
    Layout --> DocsPage
    Layout --> LoginPage
    Layout --> SignupPage
    Layout --> Footer
    Layout --> SearchModal

    HomePage --> Navbar
    HomePage --> Hero
    HomePage --> Comparison
    HomePage --> Features

    Navbar --> Button
    Hero --> Button
    DocsPage --> Sidebar
    DocsPage --> CodeBlock

    style Layout fill:#0d1117,color:#e6edf3,stroke:#30363d
    style HomePage fill:#161b22,color:#e6edf3,stroke:#30363d
    style DocsPage fill:#161b22,color:#e6edf3,stroke:#30363d
    style LoginPage fill:#161b22,color:#e6edf3,stroke:#30363d
    style SignupPage fill:#161b22,color:#e6edf3,stroke:#30363d
```

---

## 2.3 Explanation

### Package Dependencies

| Package | Depends On | Why |
|---------|-----------|-----|
| `cmd/api` | `judge`, `shared/redis`, `uuid`, `net/http`, `encoding/json` | The API server needs the data models (`judge`), the Redis wrapper to push/pop, UUID generation, and the HTTP server. |
| `cmd/worker` | `judge`, `processSubmission`, `shared/redis`, `encoding/json` | The worker deserializes submissions, processes them, and pushes results. |
| `processSubmission` | `judge`, `runBatch`, `os/exec` | The orchestrator needs data models, the batch runner, and `os/exec` to compile code. |
| `runBatch` | `judge`, `syscall` | The execution engine needs data models and `syscall.Rusage` for memory measurement. |
| `shared/redis` | `go-redis/v9` | Thin wrapper around the Redis client library. |

### Key Design Decisions

1. **`judge` is the dependency root** — It defines the data contracts and is imported by every other package. It imports nothing from the project. This is a clean "domain model" layer.

2. **`processSubmission` depends on `runBatch`, not the reverse** — The orchestrator delegates to the execution engine, creating a clear one-way dependency.

3. **`cmd/api` and `cmd/worker` are independent** — They do not import each other. They communicate exclusively through Redis queues that carry serialized `judge` structs. This enables independent scaling.

4. **The frontend is fully decoupled** — It communicates with the backend exclusively via HTTP (`/submit` and `/status`). There are no shared types or imports between the Go backend and the Next.js frontend.
