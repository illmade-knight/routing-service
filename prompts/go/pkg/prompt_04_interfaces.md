### **prompt\_04\_routing\_interfaces.md**

# **Prompt 04: Routing Service Interfaces**

**Objective**: Generate the files that define the public contract for the routing service. These interfaces form the core architectural abstraction layer.

**TDD Phase**: N/A (Foundational Interfaces)

## **Primary Prompt (Concise)**

Your task is to generate two Go files that define the public contract for the routing service:

1. **pkg/routing/types\_routing.go**: This file must contain the ConnectionInfo and DeviceToken structs.
2. **pkg/routing/interfaces\_routing.go**: This file must define the following interfaces: IngestionProducer, DeliveryProducer, PushNotifier, and MessageStore.

The methods on these interfaces must be generic and use the urn.URN and transport.SecureEnvelope types from our foundational models.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go architect designing a highly-testable, modular, and provider-agnostic microservice.

**Overall Goal**: We are creating the abstraction layer that separates our core routing logic from the specifics of any cloud provider or external service. These interfaces are the most critical part of that design.

**Your Specific Task**: Your task is to generate the complete code for two separate Go files: pkg/routing/types\_routing.go and pkg/routing/interfaces\_routing.go.

**Key Requirements to Fulfill**:

* **L2-1.2 (Provider Abstraction)**: This is the primary driver for this task. The interfaces you create are the embodiment of this requirement. They **must not** contain any types or concepts specific to a single cloud provider (e.g., no \*pubsub.Topic from the Google SDK).
* **Clarity and Purpose**:
    * IngestionProducer: Should define a method to publish an envelope into the main ingestion pipeline.
    * DeliveryProducer: Should define a method to publish an envelope to a specific, dynamic delivery topic.
    * PushNotifier: Should define a method to send a notification for an envelope to a list of device tokens.
    * MessageStore: **Must** define methods for StoreMessages, RetrieveMessages, and DeleteMessages, all keyed by a recipient's urn.URN.

**Dependencies**:

* You should assume that the urn and transport packages are available for import.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created from the existing routing-service source code.