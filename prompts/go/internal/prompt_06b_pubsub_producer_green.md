# **Prompt 06b: PubsubProducer Implementation (Green)**

**Objective**: Generate the complete implementation for the Pub/Sub Producer. The generated code must satisfy the previously defined unit tests.

**TDD Phase**: Green (Write the code to make the failing tests pass).

## **Primary Prompt (Concise)**

Your task is to generate the Go code for the internal/platform/pubsub/producer\_pubsub.go file.

I have provided the routing.IngestionProducer interface and a complete unit test file (producer\_pubsub\_test.go). Your implementation of the Producer struct **must** make all the tests in the provided test file pass.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer specializing in writing clean, testable, production-ready adapters for Google Cloud Platform services.

**Overall Goal**: We have an executable specification in the form of a failing unit test (producer\_pubsub\_test.go). Your task is to write the Go code that implements the routing.IngestionProducer interface and satisfies that specification.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **internal/platform/pubsub/producer\_pubsub.go**.

**Key Requirements to Fulfill**:

* **Pass All Tests**: This is the primary requirement.  
* **Adapter Interface for Mocking (L2-1.2)**: To enable unit testing, you **must** first define a local, unexported interface (e.g., pubsubTopicClient) that abstracts the methods used from the real \*pubsub.Topic client. The Producer struct **must** hold an instance of this interface, not the concrete \*pubsub.Topic type.  
* **Implementation Logic**: The Publish method **must** perform the following steps in order:  
  1. Marshal the input SecureEnvelope to a JSON byte slice.  
  2. Return an error if marshaling fails.  
  3. Create a pubsub.Message struct with the JSON byte slice as its Data.  
  4. Call the Publish method on the topic client.  
  5. Call Get(ctx) on the returned PublishResult to wait for completion.  
  6. Return any error from the Get call, wrapped with additional context.

**Dependencies**:

1. **The Contract**: pkg/routing/interfaces\_routing.go  
2. **The Specification**: internal/platform/pubsub/producer\_pubsub\_test.go  
3. **The Data Models**: pkg/transport/envelope.go

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created. Mandates the use of a local interface wrapper around the Pub/Sub client to ensure the implementation is unit-testable.