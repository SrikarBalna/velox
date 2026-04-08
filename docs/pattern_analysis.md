# Design Patterns & SOLID Strategy Analysis

The Velox backend uses several industry-standard design patterns to ensure the system is performant, testable, and maintainable.

---

## 1. Strategy Pattern

**Location**: `backend/processSubmission/processSubmission.go`

Velox must handle multiple programming languages, each with unique compilation and setup needs.
- **The Interface**: `LanguageStrategy` defines the contract (`Prepare`).
- **The Concrete Strategies**: `CStrategy`, `CPPStrategy`, `JavaStrategy`, `PythonStrategy`, etc., implement this interface.
- **The Context**: `SubmissionService` uses these strategies via the `StrategyRegistry` to prepare code without knowing the details of "how" it's done for each language.

> [!TIP]
> This pattern allows us to add support for a new language (e.g., Rust) simply by creating a new Strategy class, without ever touching the core submission service.

---

## 2. Adapter Pattern

**Location**: `backend/processSubmission/local_storage_adapter.go`

To maintain environment independence, the system doesn't call `os` filesystem commands directly in the business logic.
- **The Target (Port)**: `FileStorage` and `Workspace` interfaces define how the system expects to interact with storage.
- **The Adapter**: `LocalTempStorageAdapter` and `LocalWorkspace` adapt the OS's temporary filesystem to fit these interfaces.

> [!NOTE]
> This "Hexagonal" approach means we could replace the local disk with an In-Memory RAM adapter or S3 Adapter easily in the future.

---

## 3. SOLID Principles

### **S**: Single Responsibility Principle
- **Example**: `ResultAggregator` (in `processSubmission.go`) has exactly **one** job: determining the overall "Accepted" or "Wrong Answer" status from a list of test results. It does not handle database calls or file cleanup.

### **O**: Open/Closed Principle
- **Example**: `DefaultStrategyRegistry`. It is **closed** for modification (the lookup logic is fixed) but **open** for extension (you can register new languages at runtime without changing the registry class).

### **L**: Liskov Substitution Principle
- **Example**: All language implementations of `LanguageStrategy`. The `SubmissionService` can swap a `PythonStrategy` for a `JavaStrategy` at runtime, and the system behaves consistently because all adapters follow the strict interface contract.

### **I**: Interface Segregation Principle
- **Example**: Instead of one giant `Storage` interface, we have `FileStorage` (creating workspaces) and `Workspace` (interacting with files). A component that only needs to write files doesn't need to know how to create new system-level workspaces.

### **D**: Dependency Inversion Principle
- **Example**: `SubmissionService` (the high-level logic) depends on abstractions (`BatchRunner`, `StrategyRegistry`) rather than concrete worker implementations.

---

## 4. Other Key Patterns

### **Composition over Inheritance**
- **Location**: `InterpreterStrategyManager` (in `processSubmission.go`).
- **Usage**: `PythonStrategy` and `NodeStrategy` **contain** an `InterpreterStrategyManager` to share common script setup logic. This is cleaner than deep inheritance trees.

### **Repository Pattern**
- **Location**: `backend/auth/repository/`.
- **Logic**: `APIKeyRepository` and `UserRepository` wrap the SQL implementation. If you migrate from PostgreSQL to MongoDB, you only change the Repository, and your Auth Services remain untouched.

### **Middleware Pattern (Decorator)**
- **Location**: `backend/auth/middleware/`.
- **Logic**: Standard Go pattern where handlers are wrapped in decorators (`RequireAuth`, `SecurityHeaders`). Each layer adds a specific concern (Authentication, Rate Limiting, Headers) before the request hits the handler.

### **Producer-Consumer (Messaging)**
- **Location**: `shared/redis/redis.go`.
- **Logic**: Decouples the API from the Worker. The API "Produces" a job to Redis, and the Worker "Consumes" it. This allows the system to handle sudden spikes in traffic without crashing the API.
