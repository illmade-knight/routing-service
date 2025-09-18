### **prompt\_09b\_api\_handlers\_green.md**

# **Prompt 09b: API Handlers Implementation (Green)**

**Objective**: Generate the complete implementation for the API Handlers. The generated code must satisfy the previously defined unit tests.

**TDD Phase**: Green (Write the code to make the failing tests pass).

## **Primary Prompt (Concise)**

Your task is to generate the Go code for the internal/api/routing\_handlers.go file.

I have provided a complete unit test file (routing\_handlers\_test.go). Your implementation of the API struct and its handler methods (SendHandler, GetMessagesHandler) **must** make all the tests in the provided file pass.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer implementing a secure and robust HTTP API for a backend microservice.

**Overall Goal**: We have an executable specification in the form of a failing unit test (routing\_handlers\_test.go). Your task is now to write the Go code for the API handlers that satisfies that specification.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **internal/api/routing\_handlers.go**.

**Key Requirements to Fulfill**:

* **Pass All Tests**: This is the most critical requirement.  
* **SendHandler Logic**: Implement the full logic for the /send endpoint:  
  * Decode the JSON request body into a transport.SecureEnvelope.  
  * Validate that envelope.RecipientID is not zero.  
  * Publish the envelope using the ingestionProducer.  
  * Return the correct HTTP status codes for all success and error conditions.  
* **GetMessagesHandler Logic**: Implement the full logic for the /messages endpoint:  
  * Read and validate the X-User-ID header, parsing it into a urn.URN.  
  * Call messageStore.RetrieveMessages to fetch the data.  
  * Handle the case where no messages are found by returning http.StatusNoContent.  
  * If messages are found, serialize them to a JSON array in the response body.  
* **Asynchronous Cleanup (L1-2.3)**: After successfully sending the message data in the GetMessagesHandler response, you **must** launch a new, non-blocking goroutine to call messageStore.DeleteMessages. The handler must not wait for this cleanup task to complete.

**Dependencies**:

1. **The Specification**: internal/api/routing\_handlers\_test.go  
2. **The Contracts**: pkg/routing/interfaces\_routing.go.  
3. **The Data Models**: pkg/transport/envelope.go and pkg/urn/urn.go.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created. Mandates the asynchronous, fire-and-forget deletion of messages after retrieval to fulfill requirement L1-2.3.