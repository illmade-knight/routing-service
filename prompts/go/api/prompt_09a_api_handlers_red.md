### **prompt\_09a\_api\_handlers\_red.md**

# **Prompt 09a: API Handlers Unit Tests (Red)**

**Objective**: Generate a complete **unit test** file (routing\_handlers\_test.go) for the API handlers. This test will serve as the executable specification for the service's public HTTP endpoints.

**TDD Phase**: Red (Generate failing tests before the code exists).

## **Primary Prompt (Concise)**

Your task is to generate the Go code for internal/api/routing\_handlers\_test.go.

This file must contain a unit test suite for the API handlers. The tests **must** use net/http/httptest to simulate HTTP requests and **mocks for all dependencies** (IngestionProducer, MessageStore) to validate the behavior of both the SendHandler and GetMessagesHandler in isolation.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer with deep experience testing HTTP APIs using the standard library's httptest package and the testify/mock suite.

**Overall Goal**: We are now building the "glue" layer of the service. The first step is to define the precise behavior of our public-facing HTTP API. We will do this by writing a unit test that covers all success and failure paths for our handlers.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **internal/api/routing\_handlers\_test.go**.

**Key Requirements to Test**:

* **Mocking**: The test setup **must** create mock implementations for the routing.IngestionProducer and routing.MessageStore interfaces.  
* **SendHandler Test Cases**: You must test the following scenarios:  
  1. **Success**: With a valid SecureEnvelope in the request body, assert that the mock IngestionProducer.Publish is called once and the HTTP response status is http.StatusAccepted.  
  2. **Invalid JSON**: With a malformed JSON body, assert the status is http.StatusBadRequest and the producer is **never** called.  
  3. **Invalid Envelope (Missing Recipient)**: With a valid JSON body that lacks a recipientId, assert the status is http.StatusBadRequest and the producer is **never** called.  
  4. **Producer Failure**: When the mock producer's Publish method returns an error, assert the status is http.StatusInternalServerError.  
* **GetMessagesHandler Test Cases**: You must test the following scenarios:  
  1. **Success (With Messages)**: When the mock MessageStore.RetrieveMessages returns a slice of envelopes, assert the status is http.StatusOK, the JSON response body is correct, and that MessageStore.DeleteMessages is eventually called (use require.Eventually to handle its async nature).  
  2. **Success (No Messages)**: When the mock store returns an empty slice, assert the status is http.StatusNoContent and that DeleteMessages is **never** called.  
  3. **Missing Header**: When the X-User-ID header is not present, assert the status is http.StatusBadRequest and the store is **never** called.  
  4. **Store Retrieval Failure**: When RetrieveMessages returns an error, assert the status is http.StatusInternalServerError.

**Dependencies**:

* You should assume the pkg/routing, pkg/transport, and pkg/urn packages are available for import.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created. Focuses on unit testing the HTTP handlers in isolation using httptest and mocks to specify the complete API contract.