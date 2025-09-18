### **prompt\_10b\_service\_wrapper\_green.md**

# **Prompt 10b: Service Wrapper Implementation (Green)**

**Objective**: Generate the complete implementation for the routingservice.Wrapper. The generated code must satisfy the previously defined unit tests.

**TDD Phase**: Green (Write the code to make the failing tests pass).

## **Primary Prompt (Concise)**

Your task is to generate the Go code for the routingservice/routingservice.go file.

I have provided a complete unit test file (routingservice\_test.go) that defines the wrapper's behavior. Your implementation of the Wrapper struct and its lifecycle methods (New, Start, Shutdown) **must** make all the tests in the provided test file pass.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer responsible for assembling a robust microservice, ensuring all its internal components are correctly instantiated and managed.

**Overall Goal**: We have an executable specification in the form of a failing unit test (routingservice\_test.go). Your task is now to write the Go code that correctly wires together the API layer and the processing pipeline into a single, manageable service.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **routingservice/routingservice.go**.

**Key Requirements to Fulfill**:

* **Pass All Tests**: This is the most critical requirement.  
* **Constructor Logic (New)**: The New constructor **must** correctly instantiate all internal components. This includes:  
  1. Calling pipeline.NewService to create the processing service.  
  2. Calling api.NewAPI to create the HTTP handlers.  
  3. Creating and configuring an http.Server with the API handlers.  
  4. Assembling the final Wrapper struct with all its components.  
  5. It **must** correctly propagate any errors that occur during the instantiation of its internal components.  
* **Start Logic**: The Start method **must** start the background processing pipeline by calling w.processingService.Start(ctx). It **must also** start the HTTP server in a separate, non-blocking goroutine to listen for incoming requests.  
* **Shutdown Logic**: The Shutdown method **must** gracefully shut down both the apiServer and the processingService. It should attempt to shut down both components even if one fails, and it should correctly aggregate and return any errors that occur during the process.

**Dependencies**:

1. **The Specification**: routingservice/routingservice\_test.go  
2. **The Internal Packages**: internal/api and internal/pipeline.  
3. **The Contracts and Models**: All pkg/\* packages.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created. Mandates the correct wiring of internal components and the proper management of their lifecycles in separate goroutines.