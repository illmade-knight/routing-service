### **prompt\_10a\_service\_wrapper\_red.md**

# **Prompt 10a: Service Wrapper Unit Tests (Red)**

**Objective**: Generate a complete **unit test** file (routingservice\_test.go) for the Wrapper. This test will serve as the executable specification for the top-level assembly and lifecycle management of the entire service.

**TDD Phase**: Red (Generate failing tests before the code exists).

## **Primary Prompt (Concise)**

Your task is to generate the Go code for routingservice/routingservice\_test.go.

This file must contain a unit test suite for the Wrapper struct. The tests must use **mocks for all dependencies and internal components** (like the processing service) to validate that the New constructor correctly assembles all parts, and that the Start and Shutdown methods correctly orchestrate the service's lifecycle.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer testing the top-level assembly and lifecycle of a microservice, ensuring all components are wired together and managed correctly.

**Overall Goal**: We have defined the behavior of the core logic and the API layer. Now we need to specify the behavior of the main service Wrapper that brings everything together. We will define how it should be constructed and how it should manage its internal components.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **routingservice/routingservice\_test.go**.

**Key Requirements to Test**:

* **Mocking**: The test setup **must** create mock implementations for the service's direct dependencies (e.g., messagepipeline.MessageConsumer, routing.IngestionProducer, routing.Dependencies) and also for the internal components it creates, particularly the messagepipeline.StreamingService.  
* **New Constructor Test Cases**:  
  * **Success**: A test that calls New with valid mock dependencies and asserts that a non-nil Wrapper is returned without an error, verifying that the internal components were instantiated correctly.  
  * **Failure**: A test where the internal pipeline.NewService constructor is forced to fail, asserting that New propagates the error correctly.  
* **Start Method Test**:  
  * A test that calls Start on a Wrapper instance.  
  * It **must assert** that the Start method on the mock internal processingService is called exactly once.  
* **Shutdown Method Test**:  
  * A test that calls Shutdown on a Wrapper instance.  
  * It **must assert** that the Stop method on the mock internal processingService is called exactly once.  
  * It **must also assert** that the Shutdown method of the embedded http.Server is called.

**Dependencies**:

* You should assume all pkg/routing interfaces and the internal api and pipeline packages are available for use and mocking.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created. Focuses on testing the wiring and lifecycle management logic of the main service wrapper, rather than the business logic of its internal components.

---