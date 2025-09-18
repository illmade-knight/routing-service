### **prompt\_08a\_processor\_red.md**

# **Prompt 08a: RoutingProcessor Unit Tests (Red)**

**Objective**: Generate a complete **unit test** file (processor\_routing\_test.go) for the RoutingProcessor. This test will serve as the executable specification for the service's core online/offline routing logic.

**TDD Phase**: Red (Generate failing tests before the code exists).

## **Primary Prompt (Concise)**

Your task is to generate the Go code for the internal/pipeline/processor\_routing\_test.go file.

This file must contain a unit test suite for the processor function returned by NewRoutingProcessor. The tests **must use mocks for all dependencies** (e.g., PresenceCache, MessageStore, DeliveryProducer, etc.) to validate the two primary logic paths: the online user flow and the offline user flow.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer with deep experience in unit testing complex business logic with multiple mocked dependencies using the testify/mock suite.

**Overall Goal**: We are now at the heart of the routing service. We need to define the behavior of the core routing logic, which decides how to handle a message based on a user's presence. We will specify this behavior with a precise unit test before writing the implementation.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **internal/pipeline/processor\_routing\_test.go**.

**Key Requirements to Test**:

* **Mocking**: The test setup **must** create mock implementations for all interfaces defined in the routing.Dependencies struct.  
* **Test Scenarios**: You **must** create distinct test functions to validate the following scenarios:  
  1. **Online Flow**:  
     * Given: The mock PresenceCache is configured to return a successful ConnectionInfo object.  
     * The test **must assert** that deliveryProducer.Publish() is called exactly once with the correctly constructed topic name (e.g., "delivery-server123").  
     * The test **must assert** that messageStore.StoreMessages() is **never** called.  
  2. **Offline Flow (Successful Push)**:  
     * Given: The mock PresenceCache is configured to return an error (indicating the user is offline).  
     * The test **must assert** that messageStore.StoreMessages() is called first.  
     * Then, the test **must assert** that deviceTokenFetcher.Fetch() and pushNotifier.Notify() are both called.  
     * The test **must assert** that deliveryProducer.Publish() is **never** called.  
  3. **Offline Flow (Message Store Failure)**:  
     * Given: The mock PresenceCache returns an error, and the mock messageStore is configured to return an error on its StoreMessages call.  
     * The test **must assert** that the processor function returns this error.  
     * The test **must assert** that neither deviceTokenFetcher nor pushNotifier are ever called.  
  4. **Offline Flow (No Device Tokens)**:  
     * Given: The mock PresenceCache returns an error, but the deviceTokenFetcher also returns an error.  
     * The test **must assert** that messageStore.StoreMessages() is still called successfully.  
     * The test **must assert** that pushNotifier.Notify() is **never** called.  
     * The processor function **must** return a nil error, as the message was stored successfully.

**Dependencies**:

* You should assume the pkg/routing, pkg/transport, and pkg/urn packages are available for import.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created. Focuses on unit testing the two primary strategies (online vs. offline) and their specific error handling paths using a full suite of mocks.

---