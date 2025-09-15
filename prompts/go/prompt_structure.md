## **Prompt Generation Plan: routing-service**

This plan outlines the sequence of prompts we will create. The goal is to build the application from the ground up, starting with foundational data models and then iteratively building each component using a strict Test-Driven Development (TDD) cycle.

### **Phase 0: Project Context**

These are the initial, high-level prompts that set the rules for the entire project. We have already created these.

* prompt\_meta.md: The overall persona and grand strategy.
* prompt\_golang\_codestyle.md: The specific coding standards for Go.

---

### **Phase 1: The Foundation (Contracts & Data Models)**

These prompts generate the core data structures and interfaces. They are prerequisites for the TDD cycles that follow.

1. **prompt\_02\_urn\_model.md**
    * **Objective**: Generate the urn/urn.go file. This defines the strongly-typed URN identifier used throughout the entire ecosystem.
2. **prompt\_03\_envelope\_model.md**
    * **Objective**: Generate the transport/envelope.go file. This defines the SecureEnvelope, the primary data object the service routes.
3. **prompt\_04\_routing\_interfaces.md**
    * **Objective**: Generate the pkg/routing/interfaces\_routing.go file and related type definitions. This creates the core architectural contracts (MessageStore, PushNotifier, etc.) that decouple our business logic from concrete implementations.

---

### **Phase 2: The Implementation (Adapters & Pipeline Logic)**

This phase consists of strict "Red \-\> Green" TDD pairs to build the core logic and external service adapters.

4. **Firestore MessageStore Adapter**
    * **prompt\_05a\_firestore\_store\_red.md**: (Red) Generate an **integration test** (firestore\_test.go) that validates the MessageStore implementation against a Firestore emulator.
    * **prompt\_05b\_firestore\_store\_green.md**: (Green) Generate the firestore.go implementation to make the integration tests pass.
5. **Pub/Sub IngestionProducer Adapter**
    * **prompt\_06a\_pubsub\_producer\_red.md**: (Red) Generate a **unit test** (producer\_pubsub\_test.go) that validates the IngestionProducer implementation using mocks.
    * **prompt\_06b\_pubsub\_producer\_green.md**: (Green) Generate the producer\_pubsub.go implementation to make the unit tests pass.
6. **Pipeline EnvelopeTransformer**
    * **prompt\_07a\_transformer\_red.md**: (Red) Generate a **unit test** (transformer\_routing\_test.go) for the message validation and unmarshaling logic.
    * **prompt\_07b\_transformer\_green.md**: (Green) Generate the transformer\_routing.go implementation.
7. **Pipeline RoutingProcessor**
    * **prompt\_08a\_processor\_red.md**: (Red) Generate a **unit test** (processor\_routing\_test.go) for the core online/offline routing logic, mocking all external dependencies.
    * **prompt\_08b\_processor\_green.md**: (Green) Generate the processor\_routing.go implementation.

---

### **Phase 3: The Glue (API & Service Assembly)**

This final phase uses TDD pairs to build the user-facing API and the final service executable.

8. **API Handlers**
    * **prompt\_09a\_api\_handlers\_red.md**: (Red) Generate a **unit test** (routing\_handlers\_test.go) for the HTTP API handlers (/send, /messages) using httptest and mock dependencies.
    * **prompt\_09b\_api\_handlers\_green.md**: (Green) Generate the routing\_handlers.go implementation.
9. **Service Wrapper**
    * **prompt\_10a\_service\_wrapper\_red.md**: (Red) Generate a test (routingservice\_test.go) that validates the complete service assembly and its lifecycle methods (Start, Shutdown).
    * **prompt\_10b\_service\_wrapper\_green.md**: (Green) Generate the routingservice.go implementation.

---

### **Phase 4: End-to-End Validation**

This is the final verification step to ensure all components work together as intended.

10. **prompt\_11\_e2e\_test.md**
    * **Objective**: Generate the final, comprehensive end-to-end test (routing\_e2e\_test.go). This prompt will specify a test that runs the entire service with live emulators to validate the full message flow from ingestion to retrieval.

---
