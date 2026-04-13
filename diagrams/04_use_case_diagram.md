# 4. Use Case Diagram

This document identifies the primary **actor** in the Velox system and all **operations** they can perform.

---

## 4.1 System Use Case Diagram

```mermaid
flowchart LR
    User(("User"))

    subgraph "Velox System"
        direction TB

        subgraph "Authentication"
            UC1([Sign Up])
            UC2([Login])
            UC3([Logout])
        end

        subgraph "API Key Management"
            UC4([Create API Key])
            UC5([List API Keys])
            UC6([Update API Key])
            UC7([Delete API Key])
        end

        subgraph "Code Submission"
            UC8([Submit Code])
            UC9([Check Submission Status])
        end

        subgraph "Logs & Monitoring"
            UC10([View API Key Stats])
            UC11([View Dashboard])
            UC12([Health Check])
        end

        subgraph "Documentation"
            UC13([View Docs])
        end
    end

    User --> UC1
    User --> UC2
    User --> UC3
    User --> UC4
    User --> UC5
    User --> UC6
    User --> UC7
    User --> UC8
    User --> UC9
    User --> UC10
    User --> UC11
    User --> UC12
    User --> UC13

    style User fill:#ff5a00,color:#fff,stroke:#333
    style UC1 fill:#e0f2fe,stroke:#000000
    style UC2 fill:#e0f2fe,stroke:#000000
    style UC3 fill:#e0f2fe,stroke:#000000
    style UC4 fill:#fef3c7,stroke:#000000
    style UC5 fill:#fef3c7,stroke:#000000
    style UC6 fill:#fef3c7,stroke:#000000
    style UC7 fill:#fef3c7,stroke:#000000
    style UC8 fill:#dcfce7,stroke:#000000
    style UC9 fill:#dcfce7,stroke:#000000
    style UC10 fill:#fce7f3,stroke:#000000
    style UC11 fill:#fce7f3,stroke:#000000
    style UC12 fill:#fce7f3,stroke:#000000
    style UC13 fill:#f3e8ff,stroke:#000000
```

---

## 4.2 Internal System Flow (Include / Extend)

When the User triggers **Submit Code** or **Check Status**, the system internally performs the following operations:

```mermaid
flowchart TB
    UC8([Submit Code]) -->|"«include»"| V1([Validate Request])
    UC8 -->|"«include»"| V2([Generate Submission ID])
    UC8 -->|"«include»"| V3([Queue to Redis])
    UC8 -->|"«include»"| V4([Log API Request])

    V3 -.->|"async"| W1([Compile Source Code])
    W1 --> W2([Execute Test Cases])
    W2 --> W3([Measure Time & Memory])
    W2 --> W4([Enforce Resource Limits])
    W2 --> W5([Aggregate Results])
    W5 --> W6([Store Results in Redis])
    W1 --> W7([Clean Up Temp Files])
    W2 --> W7

    UC9([Check Status]) -->|"«include»"| R1([Poll Redis for Result])
    R1 -->|"found"| R2([Return Full Result])
    R1 -->|"timeout"| R3([Return Pending Status])
    UC9 -->|"«include»"| V4

    style UC8 fill:#dcfce7,stroke:#16a34a
    style UC9 fill:#dcfce7,stroke:#16a34a
    style V1 fill:#f1f5f9,stroke:#64748b
    style V2 fill:#f1f5f9,stroke:#64748b
    style V3 fill:#f1f5f9,stroke:#64748b
    style V4 fill:#f1f5f9,stroke:#64748b
    style W1 fill:#fef9c3,stroke:#ca8a04
    style W2 fill:#fef9c3,stroke:#ca8a04
    style W3 fill:#fef9c3,stroke:#ca8a04
    style W4 fill:#fef9c3,stroke:#ca8a04
    style W5 fill:#fef9c3,stroke:#ca8a04
    style W6 fill:#fef9c3,stroke:#ca8a04
    style W7 fill:#fef9c3,stroke:#ca8a04
    style R1 fill:#f1f5f9,stroke:#64748b
    style R2 fill:#f1f5f9,stroke:#64748b
    style R3 fill:#f1f5f9,stroke:#64748b
```

---

## 4.3 Detailed Use Case Descriptions

### UC1: Sign Up
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `POST /auth/signup` with JSON body |
| **Preconditions** | None |
| **Flow** | 1. User sends `name`, `email`, `password` <br/> 2. Server validates input (email format, password ≥ 8 chars) <br/> 3. Hash password and store in PostgreSQL <br/> 4. Return `201 Created` with user details |
| **Error Cases** | Missing name → 400, Email taken → 400, Invalid email → 400, Password too short → 400 |

### UC2: Login
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `POST /auth/login` with JSON body |
| **Flow** | 1. User sends `email` and `password` <br/> 2. Server verifies credentials <br/> 3. Generate JWT token <br/> 4. Return `200 OK` with token |
| **Error Cases** | Invalid credentials → 401 |

### UC3: Logout
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `POST /auth/logout` |
| **Flow** | 1. Stateless JWT — server responds with success <br/> 2. Client discards token locally |

### UC4: Create API Key
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `POST /auth/api-keys` (requires session auth) |
| **Flow** | 1. User provides `name`, optional `scopes` and `expires_at` <br/> 2. Default scopes: `["submit", "status"]` <br/> 3. Server generates key, stores hash in DB <br/> 4. Returns full key (shown only once), ID, and display hint |
| **Error Cases** | Missing name → 400, Unauthorized → 401 |

### UC5: List API Keys
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `GET /auth/api-keys` (requires session auth) |
| **Flow** | Returns all API keys belonging to the authenticated user |

### UC6: Update API Key
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `PATCH /auth/api-keys?id=<uuid>` (requires session auth) |
| **Flow** | Renames an existing API key |

### UC7: Delete API Key
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `DELETE /auth/api-keys?id=<uuid>` (requires session auth) |
| **Flow** | Permanently revokes and deletes an API key |

### UC8: Submit Code
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `POST /submit` with JSON body (requires API key auth) |
| **Preconditions** | Request body contains `language`, `source_code`, and at least one `test_case` |
| **Flow** | 1. API validates `TimeLimitMs ≤ 5000` and `MemoryLimitKb ≤ 512000` <br/> 2. Generate UUID via `uuid.New()` <br/> 3. Serialize request to JSON <br/> 4. `LPUSH` to Redis `"submissions"` queue <br/> 5. Log API request <br/> 6. Return `202 Accepted` with `submission_id` |
| **Postconditions** | Submission is queued for processing |
| **Error Cases** | Invalid JSON → 400, Limits too high → 400, Redis push failure → 500 |

### UC9: Check Submission Status
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `GET /status?submission_id=<id>` (requires API key auth) |
| **Flow** | 1. Extract `submission_id` from query params <br/> 2. `BRPOP "results:<id>"` with 1s timeout <br/> 3a. If found → return the full response JSON and update log <br/> 3b. If timeout → return `{"status": "pending"}` |
| **Error Cases** | Missing `submission_id` → 400 |

### UC10: View API Key Stats
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `GET /auth/api-keys/stats?id=<uuid>` (requires session auth) |
| **Flow** | Returns usage metrics (Total requests, RPM, RPD, Success Rate) for a specific API key |

### UC11: View Dashboard
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `GET /dashboard` (requires session auth) |
| **Flow** | Returns user-specific profile and activity data |

### UC12: Health Check
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `GET /health` |
| **Flow** | Returns `{"status": "healthy"}` if API server is running |

### UC13: View Swagger Docs
| Field | Value |
|-------|-------|
| **Actor** | User |
| **Trigger** | `GET /swagger/` (development mode only) |
| **Flow** | Serves the interactive Swagger UI for API exploration |

---

## 4.4 Supported Languages Matrix

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
