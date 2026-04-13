# 4. Use Case Diagram

This document identifies all **actors** in the Velox system and the **operations** they perform or participate in.

---

## 4.1 System Use Case Diagram

```mermaid
flowchart LR
    User((👤 User<br/>Developer / Student))
    Admin((🔧 Admin<br/>DevOps))
    CI((🤖 CI Pipeline<br/>GitHub Actions))

    subgraph "Velox System"
        direction TB
        subgraph "API Server"
            UC1([Submit Code])
            UC2([Check Submission Status])
            UC3([Validate Request])
            UC4([Generate Submission ID])
        end

        subgraph "Worker"
            UC5([Compile Source Code])
            UC6([Execute Test Cases])
            UC7([Measure Time & Memory])
            UC8([Enforce Resource Limits])
            UC9([Aggregate Results])
            UC10([Clean Up Temp Files])
        end

        subgraph "Redis"
            UC11([Queue Submission])
            UC12([Store Results])
            UC13([Retrieve Results])
        end

        subgraph "Frontend"
            UC14([View Landing Page])
            UC15([Read Documentation])
            UC16([Sign In])
            UC17([Sign Up])
            UC18([Search Documentation])
        end

        subgraph "Infrastructure"
            UC19([Build Docker Images])
            UC20([Run Load Tests])
            UC21([Deploy Stack])
            UC22([Run CI Tests])
        end
    end

    User -->|"POST /submit"| UC1
    User -->|"GET /status"| UC2
    User --> UC14
    User --> UC15
    User --> UC16
    User --> UC17
    User --> UC18

    UC1 --> UC3
    UC1 --> UC4
    UC1 --> UC11

    UC2 --> UC13

    UC11 -.->|"async"| UC5
    UC5 --> UC6
    UC6 --> UC7
    UC6 --> UC8
    UC6 --> UC9
    UC9 --> UC12
    UC5 --> UC10
    UC6 --> UC10

    Admin --> UC19
    Admin --> UC20
    Admin --> UC21

    CI --> UC22

    style User fill:#ff5a00,color:#fff,stroke:#333
    style Admin fill:#2563eb,color:#fff,stroke:#333
    style CI fill:#16a34a,color:#fff,stroke:#333
```

---

## 4.2 Detailed Use Case Descriptions

### UC1: Submit Code
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `POST /submit` with JSON body |
| **Preconditions** | Request body contains `language`, `source_code`, and at least one `test_case` |
| **Flow** | 1. API validates `TimeLimitMs ≤ 5000` and `MemoryLimitKb ≤ 512000` <br/> 2. Generate UUID via `uuid.New()` <br/> 3. Serialize request to JSON <br/> 4. `LPUSH` to Redis `"submissions"` queue <br/> 5. Return `202 Accepted` with `submission_id` |
| **Postconditions** | Submission is queued for processing |
| **Error Cases** | Invalid JSON → 400, Limits too high → 400, Redis push failure → 500 |

### UC2: Check Submission Status
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `GET /status?submission_id=<id>` |
| **Preconditions** | A submission with the given ID was previously submitted |
| **Flow** | 1. Extract `submission_id` from query params <br/> 2. `BRPOP "results:<id>"` with 1s timeout <br/> 3a. If found → return the full response JSON <br/> 3b. If timeout → return `{"status": "pending"}` |
| **Postconditions** | Client receives result or pending status |
| **Error Cases** | Missing `submission_id` → 400 |

### UC3: Validate Request
| Field | Value |
|-------|-------|
| **Actor** | API Server (internal) |
| **Checks** | HTTP method is POST, JSON is valid, `TimeLimitMs ≤ 5000`, `MemoryLimitKb ≤ 512000` |

### UC5: Compile Source Code
| Field | Value |
|-------|-------|
| **Actor** | Worker (internal) |
| **Supported Languages** | C (gcc), C++ (g++), Java (javac), TypeScript (tsc), C# (dotnet build) |
| **Interpreted Languages** | Python, Node.js — no compilation, file is written to temp dir and executed directly |
| **Error Output** | Compiler stderr is captured and returned as `CompileError` field |

### UC6: Execute Test Cases
| Field | Value |
|-------|-------|
| **Actor** | Worker → RunBatch (internal) |
| **Flow** | For each test case: <br/> 1. Create context with timeout <br/> 2. Run binary/script with `stdin` piped from test input <br/> 3. Capture `stdout`, `stderr` <br/> 4. Measure CPU time and peak memory <br/> 5. Compare actual output to expected output |
| **Possible Statuses** | `Accepted`, `Wrong Answer`, `Runtime Error`, `Time Limit Exceeded`, `Memory Limit Exceeded` |

### UC8: Enforce Resource Limits
| Field | Value |
|-------|-------|
| **Actor** | RunBatch (internal) |
| **Time Limit** | Enforced via `context.WithTimeout`. Default: 3000ms. Max: 5000ms. |
| **Memory Limit** | Measured via `syscall.Rusage.Maxrss`. Default: 256MB. Max: 512MB. Platform-aware (macOS divides by 1024). |

### UC20: Run Load Tests
| Field | Value |
|-------|-------|
| **Actor** | Admin / Developer |
| **Tool** | `tests/load_test_cmd.go` |
| **Behavior** | Sends 20 concurrent submissions per language (7 languages × 20 = 140 total), polls for results, and prints percentile-based performance metrics (P50, P90, P95, P99). |

### UC22: Run CI Tests
| Field | Value |
|-------|-------|
| **Actor** | GitHub Actions CI pipeline |
| **Trigger** | Push to `main` or Pull Request to `main` |
| **Strategy** | 3-job pipeline: (1) Matrix unit tests per language (C, CPP, Java, Node, Python, TS, CSharp), (2) Build load test binary from tests/load/, (3) Build performance test binary from tests/performance/ |
| **Command** | `go test -v -race ./processSubmission -run TestProcessSubmission_<Language>$` |

---

## 4.3 Supported Languages Matrix

The system supports 7 programming languages. Each language follows a specific execution pipeline:

```mermaid
graph LR
    subgraph "Compiled Languages"
        C["C<br/>gcc → binary"]
        CPP["C++<br/>g++ → binary"]
        Java["Java<br/>javac → java -cp"]
        TS["TypeScript<br/>tsc → node"]
        CS["C#<br/>dotnet build → dotnet run"]
    end

    subgraph "Interpreted Languages"
        Python["Python<br/>python3 script.py"]
        Node["Node.js<br/>node script.js"]
    end

    C --> RB["RunBatch"]
    CPP --> RB
    Java --> RB
    TS --> RB
    CS --> RB
    Python --> RB
    Node --> RB

    style C fill:#555555,color:#fff
    style CPP fill:#004482,color:#fff
    style Java fill:#f89820,color:#000
    style TS fill:#3178c6,color:#fff
    style CS fill:#68217a,color:#fff
    style Python fill:#3776ab,color:#fff
    style Node fill:#339933,color:#fff
    style RB fill:#ff5a00,color:#fff
```
