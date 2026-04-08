# 6. Component & Deployment Diagram

This document shows the **physical deployment topology** — Docker containers, networks, ports, and how the system is built and deployed.

---

## 6.1 Docker Deployment Diagram

```mermaid
graph TB
    subgraph "Host Machine"
        subgraph "Docker Compose Network (bridge)"
            
            subgraph "redis (redis:alpine)"
                R_PROC["redis-server<br/>--requirepass 123456"]
                R_PORT["Port 6379 internal<br/>→ 6380 host"]
                R_HEALTH["Healthcheck:<br/>redis-cli -a 123456 ping<br/>interval: 5s, retries: 5"]
            end

            subgraph "api (Dockerfile.api)"
                A_BUILD["Build Stage: golang:latest<br/>CGO_ENABLED=0 GOOS=linux<br/>go build ./cmd/api"]
                A_RUNTIME["Runtime: debian:bookworm-slim<br/>User: apiuser<br/>Binary: /app/main"]
                A_PORT["Port 8080 internal<br/>→ 8080 host"]
                A_ROUTES["Routes:<br/>POST /submit<br/>GET /status"]
            end

            subgraph "velox (Dockerfile.worker)"
                W_BUILD["Build Stage: golang:latest<br/>CGO_ENABLED=0 GOOS=linux<br/>go build ./cmd/worker"]
                W_RUNTIME["Runtime: debian:bookworm-slim<br/>User: runner<br/>Binary: /app/main"]
                W_TOOLS["Installed Compilers/Runtimes:<br/>• gcc, g++<br/>• OpenJDK 17<br/>• Python 3<br/>• Node.js + npm<br/>• .NET SDK 8.0<br/>• TypeScript (tsc)"]
            end
        end

        Client["👤 Client<br/>localhost:8080"]
    end

    Client -->|"HTTP"| A_PORT
    A_RUNTIME -->|"LPUSH / BRPOP"| R_PROC
    W_RUNTIME -->|"BRPOP / LPUSH"| R_PROC

    A_RUNTIME -.->|"depends_on"| R_HEALTH
    W_RUNTIME -.->|"depends_on"| R_HEALTH

    style R_PROC fill:#dc2626,color:#fff
    style A_RUNTIME fill:#ff5a00,color:#fff
    style W_RUNTIME fill:#2563eb,color:#fff
    style Client fill:#16a34a,color:#fff
```

---

## 6.2 Build Pipeline Diagram

Shows how the Docker images are constructed with multi-stage builds.

```mermaid
graph LR
    subgraph "Dockerfile.api (Multi-Stage)"
        A1["Stage 1: golang:latest<br/><b>Builder</b>"]
        A2["COPY backend/go.mod<br/>go mod download"]
        A3["COPY backend/<br/>go build ./cmd/api"]
        A4["Stage 2: debian:bookworm-slim<br/><b>Runtime</b>"]
        A5["COPY --from=builder /build/backend/main"]
        A6["USER apiuser<br/>EXPOSE 8080<br/>ENTRYPOINT main"]
    end

    A1 --> A2 --> A3 --> A4 --> A5 --> A6

    subgraph "Dockerfile.worker (Multi-Stage)"
        W0["Stage 0: mcr.microsoft.com/dotnet/sdk:8.0<br/><b>.NET SDK source</b>"]
        W1["Stage 1: golang:latest<br/><b>Builder</b>"]
        W2["COPY backend/<br/>go build ./cmd/worker"]
        W3["Stage 2: debian:bookworm-slim<br/><b>Runtime</b>"]
        W4["apt-get install:<br/>gcc g++ python3 nodejs npm openjdk-17"]
        W5["COPY --from=dotnet-sdk<br/>/usr/share/dotnet"]
        W6["npm install -g typescript"]
        W7["COPY --from=builder<br/>/build/backend/main"]
        W8["USER runner<br/>ENTRYPOINT main"]
    end

    W0 --> W5
    W1 --> W2 --> W3 --> W4 --> W5 --> W6 --> W7 --> W8

    style A1 fill:#00add8,color:#fff
    style A4 fill:#ff5a00,color:#fff
    style W0 fill:#68217a,color:#fff
    style W1 fill:#00add8,color:#fff
    style W3 fill:#2563eb,color:#fff
```

---

## 6.3 CI/CD Pipeline Diagram

```mermaid
graph TD
    subgraph "GitHub"
        Push["Push to main /<br/>Pull Request to main"]
    end

    subgraph "GitHub Actions Runner (ubuntu-latest)"
        Checkout["actions/checkout@v4"]
        SetupGo["actions/setup-go@v5<br/>go 1.24.x"]
        SetupDotnet["actions/setup-dotnet@v4<br/>.NET 8.0.x"]
        InstallTS["npm install -g typescript"]

        subgraph "Job 1: Unit Tests — Matrix Strategy (fail-fast: false)"
            TestC["go test -race<br/>TestProcessSubmission_C"]
            TestCPP["go test -race<br/>TestProcessSubmission_CPP"]
            TestJava["go test -race<br/>TestProcessSubmission_Java"]
            TestNode["go test -race<br/>TestProcessSubmission_Node"]
            TestPython["go test -race<br/>TestProcessSubmission_Python"]
            TestTS["go test -race<br/>TestProcessSubmission_TS"]
            TestCSharp["go test -race<br/>TestProcessSubmission_CSharp"]
        end

        subgraph "Job 2: Build Load Tests (needs: unit-tests)"
            LoadBuild["go build -o load_test_binary<br/>tests/load/"]
            LoadArtifact["Upload Artifact:<br/>load-test-binary"]
        end

        subgraph "Job 3: Build Performance Tests (needs: unit-tests)"
            PerfBuild["go build -o performance_test_binary<br/>tests/performance/"]
            PerfArtifact["Upload Artifact:<br/>performance-test-binary"]
        end
    end

    Push --> Checkout --> SetupGo --> SetupDotnet --> InstallTS

    InstallTS --> TestC
    InstallTS --> TestCPP
    InstallTS --> TestJava
    InstallTS --> TestNode
    InstallTS --> TestPython
    InstallTS --> TestTS
    InstallTS --> TestCSharp

    TestC --> LoadBuild
    TestCPP --> LoadBuild
    TestJava --> LoadBuild
    TestNode --> LoadBuild
    TestPython --> LoadBuild
    TestTS --> LoadBuild
    TestCSharp --> LoadBuild

    TestC --> PerfBuild
    TestCPP --> PerfBuild
    TestJava --> PerfBuild
    TestNode --> PerfBuild
    TestPython --> PerfBuild
    TestTS --> PerfBuild
    TestCSharp --> PerfBuild

    LoadBuild --> LoadArtifact
    PerfBuild --> PerfArtifact

    style Push fill:#333,color:#fff
    style TestC fill:#555,color:#fff
    style TestCPP fill:#004482,color:#fff
    style TestJava fill:#f89820,color:#000
    style TestNode fill:#339933,color:#fff
    style TestPython fill:#3776ab,color:#fff
    style TestTS fill:#3178c6,color:#fff
    style TestCSharp fill:#68217a,color:#fff
    style LoadBuild fill:#16a34a,color:#fff
    style LoadArtifact fill:#16a34a,color:#fff
    style PerfBuild fill:#f59e0b,color:#000
    style PerfArtifact fill:#f59e0b,color:#000
```

---

## 6.4 Full System Architecture Diagram

This combines the frontend, backend, and infrastructure into one view.

```mermaid
graph TB
    subgraph "User's Browser"
        FE["Next.js Frontend<br/>(React, TailwindCSS)<br/>Pages: Home, Docs, Login, Signup"]
    end

    subgraph "Docker Compose Stack"
        subgraph "API Service (port 8080)"
            API["Go HTTP Server<br/>cmd/api/main.go<br/>Endpoints: /submit, /status, /health"]
            AUTH["Auth Module<br/>JWT + API Key Middleware<br/>Signup, Login, Dashboard"]
        end

        subgraph "PostgreSQL"
            PG["Neon PostgreSQL<br/>Users + API Keys"]
        end

        subgraph "Redis Service (port 6379)"
            REDIS["Redis Alpine<br/>Password: protected<br/>Pool: 1000 max, 100 idle"]
            SQ["Queue: 'submissions'"]
            RQ["Queue: 'results:&lt;id&gt;'"]
        end

        subgraph "Worker Service"
            WORKER["Go Worker Process<br/>cmd/worker/main.go<br/>Infinite polling loop"]
            PS["processSubmission<br/>Strategy Pattern routing + compilation"]
            RB["runBatch<br/>Execution + measurement"]
            COMPILERS["gcc | g++ | javac | dotnet<br/>python3 | node | tsc"]
        end
    end

    subgraph "Testing"
        LT["tests/load/main.go<br/>Concurrent load testing<br/>7 languages × 2 each"]
        PT["tests/performance/main.go<br/>Sequential performance benchmarking<br/>Detailed per-language metrics"]
    end

    FE -->|"HTTP POST/GET"| API
    LT -->|"HTTP POST/GET"| API
    PT -->|"HTTP POST/GET"| API
    API -->|"LPUSH"| SQ
    API -->|"BRPOP"| RQ
    AUTH -->|"SQL"| PG
    WORKER -->|"BRPOP"| SQ
    WORKER --> PS
    PS --> RB
    RB --> COMPILERS
    WORKER -->|"LPUSH"| RQ

    style FE fill:#0d1117,color:#e6edf3,stroke:#30363d
    style API fill:#ff5a00,color:#fff
    style AUTH fill:#7c3aed,color:#fff
    style PG fill:#336791,color:#fff
    style REDIS fill:#dc2626,color:#fff
    style WORKER fill:#2563eb,color:#fff
    style LT fill:#16a34a,color:#fff
    style PT fill:#f59e0b,color:#000
```

---

## 6.5 Environment Variables

| Variable | Service | Default | Description |
|----------|---------|---------|-------------|
| `REDIS_ADDR` | api, worker | `localhost:6379` | Redis connection address |
| `REDIS_PASSWORD` | api, worker | `123456` | Redis authentication password |

---

## 6.6 Port Mapping

| Service | Internal Port | External Port | Protocol |
|---------|--------------|---------------|----------|
| Redis | 6379 | 6380 | TCP |
| API Server | 8080 | 8080 | HTTP |
| Worker | — | — | No exposed ports |

---

## 6.7 Explanation

### Why This Architecture?

| Decision | Rationale |
|----------|-----------|
| **Separate API and Worker containers** | The API is lightweight (just HTTP + Redis). The Worker is heavy (compilers + runtimes). Separating them allows independent scaling — you can run 1 API and 10 Workers behind the same Redis queue. |
| **Redis as the message queue** | Simple, battle-tested, and already provides `BRPOP` for blocking pops. No need for Kafka or RabbitMQ at this scale. |
| **Multi-stage Docker builds** | The final API image is ~50 MB (just a static Go binary + CA certs). The worker image is larger due to compilers. Multi-stage builds keep the image sizes minimal. |
| **Non-root users** | Both containers run as non-root (`apiuser` / `runner`) for security. This is especially important for the worker, which executes untrusted user code. |
| **Healthcheck on Redis** | Docker Compose uses healthchecks to ensure Redis is ready before starting the API and Worker. This prevents race conditions on startup. |
| **No persistent storage** | All data is transient in Redis. After results are consumed, they are gone. This is intentional — Velox is a stateless execution engine, not a database. If you need persistence, add a database layer on top. |
