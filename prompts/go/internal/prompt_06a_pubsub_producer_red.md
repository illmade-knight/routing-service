# **Prompt 06a: PubsubProducer Unit Tests (Red)**

**Objective**: Generate a complete **unit test** file (producer\_pubsub\_test.go) for the Producer. This test will serve as the executable specification for our Pub/Sub adapter's implementation.

**TDD Phase**: Red (Generate failing tests before the code exists).

## **Primary Prompt (Concise)**

Your task is to generate the Go code for the internal/platform/pubsub/producer\_pubsub\_test.go file.

This file must contain a unit test suite for a Producer struct that will implement the routing.IngestionProducer interface. The tests **must** use a **mock Pub/Sub topic client** to validate that the Publish method correctly serializes a SecureEnvelope to JSON and calls the topic's Publish method with the resulting bytes, without making any real network calls.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer specializing in writing thorough unit tests using the testify/mock and testify/require libraries.

**Overall Goal**: We are continuing the TDD implementation of the routing service. The next component is the concrete Pub/Sub adapter that implements the IngestionProducer interface. We will define its behavior by writing a comprehensive unit test first.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **internal/platform/pubsub/producer\_pubsub\_test.go**.

**Key Requirements to Test**:

* **Mocking**: The test **must** define a local mock struct (e.g., mockPubsubTopicClient) that implements an interface for the Pub/Sub topic client. This mock will be used to assert behavior and control outcomes.  
* **Test Cases (Table-Driven)**: You **must** use a table-driven test to cover the following scenarios for the Publish method:  
  1. **Success Case**: The test must verify that:  
     * The input SecureEnvelope is correctly marshaled to JSON.  
     * The mock topic's Publish method is called exactly once with a pubsub.Message containing the correct JSON payload.  
     * The Publish method returns no error.  
  2. **JSON Marshal Failure**: The test must simulate a JSON marshaling error. It **must assert** that an error is returned and that the mock topic's Publish method is **never** called.  
  3. **Pub/Sub Publish Failure**: The test must simulate the Pub/Sub result.Get() call returning an error. It **must assert** that this error is correctly propagated back to the caller.

**Dependencies**:

* You should assume the pkg/routing, pkg/transport, and pkg/urn packages, containing the interfaces and data models, are available for import.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created. Focuses on pure unit testing using a mock client to ensure the adapter's serialization and error handling logic is correct.

---