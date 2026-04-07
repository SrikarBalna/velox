# Velox

Velox is a high-performance, containerized code execution engine (Online Judge) built with Go and Docker. It allows you to submit code in various languages, execute it against multiple test cases, and receive detailed resource usage (time and memory) along with execution status.

## 📁 Project Structure

The project is split into two main services: an **API Server** and a **Worker**.

### `backend/`

Core Go application logic.

- **`cmd/`**
  - **`api/`**: Entry point for the HTTP server. Handles `/submit` and `/status` endpoints.
  - **`worker/`**: Entry point for the background worker. Continuously polls Redis for new submissions.
- **`judge/`**: Defines the data models (`judje.go`) for requests and responses used across the system.
- **`processSubmission/`**: The language orchestrator. Handles compilation (for C, C#, C++, Java, TS) and script preparation (for Python, Node).
- **`runBatch/`**: The execution engine. Runs binaries/scripts in a controlled environment, pipes input, and captures results, Time (ms), and Memory (KB).
- **`shared/redis/`**: Basic Redis wrapper for pushing/popping submissions and results.
- **`docs/`**: Automatically generated Swagger (OpenAPI 3.0) documentation.

---

## 📖 API Documentation

The project uses Swagger (OpenAPI 3.0) for API documentation.

### 🔍 View Documentation
When the server is running in development mode (`GO_ENV=development`), you can access the interactive Swagger UI at:
- **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

### 🛠 Generate Documentation
If you add new endpoints or update annotations, regenerate the documentation using:
```bash
cd backend
go run github.com/swaggo/swag/cmd/swag init -g cmd/api/main.go
```

### `build/`

Contains infrastructure configuration.

- **`Dockerfile.api`**: Minimal Go runtime for the API service.
- **`Dockerfile.worker`**: Full-featured image containing compilers (gcc, g++, javac) and runtimes (python, node, openjdk) required for judging.

### `docker-compose.yml`

Orchestrates the `api`, `worker`, and `redis` services into a single local environment.

---

## 🛠 Technology Stack

- **Language**: Go (Golang)
- **Queue/Store**: Redis
- **Containerization**: Docker & Docker Compose
- **Supported Languages**:
  - C (GCC 12+)
  - C++ (G++ 12+)
  - Java (OpenJDK 17)
  - C#
  - Python (3.x)
  - Node.js
  - TypeScript

---

## 🚀 Quick Start

Ensure you have **Docker** and **Docker Compose** installed.

1. **Clone the repository:**

   ```bash
   git clone https://github.com/RISHIK92/velox.git
   cd velox
   ```

2. **Spin up the stack:**

   ```bash
   docker compose up --build
   ```

3. **Test a submission:**

   ```bash
   curl -X POST http://localhost:8080/submit \
     -H "Content-Type: application/json" \
     -d '{
       "language": "cpp",
       "source_code": "#include <iostream>\nusing namespace std;\nint main() { int a, b; while (cin >> a >> b) cout << a + b << endl; return 0; }",
       "test_cases": [{"test_case_id": 1, "input": "5 10", "expected_output": "15"}]
     }'
   ```

4. **Check status:**
   ```bash
   curl http://localhost:8080/status?submission_id=<ID_FROM_PREVIOUS_STEP>
   ```

---

## 🤝 Contributing

We welcome contributions! To help you get started:

### Adding a New Language Support

1. **Update `backend/judge/judje.go`**: Add the language string to comments if needed.
2. **Update `backend/processSubmission/ProcessSubmission.go`**:
   - Add a new `case` in the `switch` block.
   - Implement a compiler/preparer function (similar to `CompileInMemoryCPP`).
3. **Update `build/Dockerfile.worker`**: Ensure the necessary compiler or runtime is installed in the worker image.

### Development Workflow

- **Code Style**: We follow standard Go formatting (`go fmt`).
- **Testing**: Use the provided `curl` samples to verify your changes.
- **Watch Mode**: You can use `docker compose watch` (if supported by your version) to automatically rebuild services when files in `backend/` change.

### Steps to Contribute

1. Fork the repo.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a Pull Request.
