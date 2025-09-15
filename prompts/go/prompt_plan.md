## **Prompt Plan: Routing Service**

### **1\. Guiding Philosophy**

Our strategy is to generate the routing-service using a **Test-Driven Development (TDD)** inspired workflow. This ensures that the generated code is correct by design and provides a clear, maintainable link between requirements and implementation.

We will use a series of smaller, self-contained prompts, each responsible for generating a single file or a closely related set of interfaces. This modular approach enhances clarity, context, testability, and maintainability.

### **2\. The Prompting Workflow: A Phased TDD Approach**

We will generate the repository in a logical order that mirrors a typical development process, ensuring that dependencies are created before they are needed. Each component within these phases will be generated using the three-step TDD cycle.

#### **Phase 1: The Foundation (Contracts & Data Models)**

* **Goal**: Generate the core data structures and interfaces that define the contracts for the entire service.
* **Modules**:
    * urn/urn.go: The URN type definition.
    * transport/envelope.go: The SecureEnvelope struct.
    * routing/types.go: Core domain types like ConnectionInfo.
    * routing/interfaces.go: The key service interfaces (IngestionProducer, MessageStore, etc.).

#### **Phase 2: The Implementation (Adapters & Pipeline Logic)**

* **Goal**: Implement the core business logic and concrete adapters for external services.
* **Modules**:
    * internal/platform/persistence/firestore.go: The Firestore implementation of MessageStore.
    * internal/platform/pubsub/producer\_pubsub.go: The Pub/Sub implementation of IngestionProducer.
    * internal/pipeline/transformer\_routing.go: The message validation and transformation logic.
    * internal/pipeline/processor\_routing.go: The core online/offline routing logic.

#### **Phase 3: The Glue (API & Service Assembly)**

* **Goal**: Generate the top-level API handlers and the service wrapper that ties all modules together.
* **Modules**:
    * internal/api/routing\_handlers.go: The HTTP handlers for /send and /messages.
    * routingservice/routingservice.go: The main service wrapper that manages startup, shutdown, and component assembly.

#### **The Core TDD Cycle (Applied in Each Phase)**

For each file generated, we will follow this strict cycle:

1. **Generate the Interface (The Contract)**: If applicable, a prompt to generate the clean Go interface for the component.
2. **Generate the Tests (The Specification)**: A prompt that takes the interface and requirements as input to generate a complete \_test.go file. This is the **"Red"** phase.
3. **Generate the Implementation (The Code)**: A prompt that takes both the interface and the failing tests as input. The instruction is simple: "Write the Go code that implements the interface and makes all the provided tests pass." This is the **"Green"** phase.

### **3\. Requirements Traceability**

To ensure prompts and code stay aligned with requirements, we will use a traceability system. Each prompt will include a metadata block that lists the specific requirement IDs (e.g., L1-1.1, L2-2.3) it is designed to fulfill. A simple script can then be used in CI to verify that all "shall" or "must" requirements are covered by at least one prompt, creating a feedback loop that keeps the system in sync.

---

