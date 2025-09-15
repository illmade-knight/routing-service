### **L2: Architectural & Non-Functional Requirements (L2\_Architectural\_Requirements.md)**

This document defines *how* the routing service must be built, outlining the rules and qualities that govern its design.

# **L2: Architectural & Non-Functional Requirements**

## **1\. Core Architecture Principles**

| ID | Requirement | Description |
| :---- | :---- | :---- |
| **1.1** | **Modularity & Separation of Concerns** | The system's core domains—API handling, message processing, and persistence—**shall** be implemented as distinct, loosely-coupled modules. |
| **1.2** | **Provider Abstraction** | The core business logic **shall** be decoupled from any specific cloud provider's SDKs or APIs. This **must** be achieved through a provider/adapter design pattern, with all external dependencies defined by interfaces (e.g., MessageStore, IngestionProducer). |
| **1.3** | **Asynchronous Pipeline Processing** | The core message routing logic **shall** be implemented as a multi-stage, asynchronous message processing pipeline (e.g., consumer \-\> transformer \-\> processor) to ensure high throughput and separation of concerns. |
| **1.4** | **Concurrency** | The system **shall** be designed to process multiple messages concurrently through the use of a configurable number of pipeline workers. |

## **2\. Reliability & Safety**

| ID | Requirement | Description |
| :---- | :---- | :---- |
| **2.1** | **Atomic Writes** | All operations that write multiple message documents to the persistence layer **shall** be performed in a single atomic transaction to prevent partial state updates. The same atomicity **must** be applied to batch deletion operations. |
| **2.2** | **Input Validation** | The system **shall** perform strict validation on all incoming messages at the earliest possible stage (the "transformer" stage). Messages that are malformed or missing critical information (like SenderID or RecipientID) **must** be rejected and prevented from entering the core processing logic. |
| **2.3** | **Graceful Shutdown** | The application **shall** support graceful shutdown. Upon receiving a termination signal (e.g., SIGINT, SIGTERM), it **must** attempt to finish processing in-flight messages and close all external connections cleanly. |

## **3\. Testability**

| ID | Requirement | Description |
| :---- | :---- | :---- |
| **3.1** | **Unit Testability** | All core logic modules **must** be unit-testable in isolation from external services, a direct result of the Provider Abstraction requirement (1.2). |
| **3.2** | **High-Fidelity Integration Testing** | The system **shall** be designed to be testable against real cloud emulators (e.g., for Pub/Sub and Firestore) to validate the full lifecycle of operations in a production-like environment. |

---