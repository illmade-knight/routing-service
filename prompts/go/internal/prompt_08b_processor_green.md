# **Prompt 08b: RoutingProcessor Implementation (Green)**

**Objective**: Generate the complete implementation for the RoutingProcessor. The generated code must satisfy the previously defined unit tests.

**TDD Phase**: Green (Write the code to make the failing tests pass).

## **Primary Prompt (Concise)**

Your task is to generate the Go code for the internal/pipeline/processor\_routing.go file.

I have provided a complete unit test file (processor\_routing\_test.go) that defines the processor's behavior. Your implementation of the NewRoutingProcessor factory function and the processor it returns **must** make all tests in the provided test file pass.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer implementing the core business logic for a real-time messaging microservice.

**Overall Goal**: We have an executable specification in the form of a failing unit test (processor\_routing\_test.go). Your task is now to write the Go code that implements the service's core routing logic to satisfy that specification.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **internal/pipeline/processor\_routing.go**.

**Key Requirements to Fulfill**:

* **Pass All Tests**: This is the most critical requirement.  
* **Online Strategy (L1-1.1)**: The processor logic **must** first attempt to fetch the user's connection info from the presenceCache. If this call succeeds (returns a nil error), the processor **must** immediately construct the delivery topic string as "delivery-" \+ connInfo.ServerInstanceID and publish the message using the deliveryProducer. It should then stop processing for that message.  
* **Offline Strategy (L1-1.2)**: If the presenceCache.Fetch call returns an error, the processor **must** execute the offline strategy.  
* **Offline Order of Operations (L2.1, L2.4)**: The offline strategy **must** follow this exact sequence:  
  1. Call messageStore.StoreMessages(). If this fails, the processor **must** return the error immediately.  
  2. Call deviceTokenFetcher.Fetch().  
  3. If device tokens are retrieved successfully, call pushNotifier.Notify().  
* **Offline Error Handling (L2.4)**: A failure from either deviceTokenFetcher.Fetch or pushNotifier.Notify **must not** result in a returned error from the processor function. The primary requirement is that the message was successfully stored; a failed push notification should not cause the entire message to be re-processed.

**Dependencies**:

1. **The Specification**: internal/pipeline/processor\_routing\_test.go  
2. **The Contracts**: pkg/routing/interfaces\_routing.go and its dependency interfaces.  
3. **The Data Models**: pkg/transport/envelope.go and pkg/urn/urn.go.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created. Mandates the specific "online vs. offline" logic flow and the critical error handling behavior for the offline path.