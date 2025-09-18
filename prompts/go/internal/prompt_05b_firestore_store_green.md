### **prompt\_05b\_firestore\_store\_green.md**

# **Prompt 05b: FirestoreStore Implementation (Green)**

**Objective**: Generate the complete implementation for the FirestoreStore. The generated code must satisfy the previously defined integration tests.

**TDD Phase**: Green (Write the code to make the failing tests pass).

## **Primary Prompt (Concise)**

Your task is to generate the Go code for the internal/platform/persistence/firestore.go file.

I have provided the routing.MessageStore interface and a complete integration test file (firestore\_test.go). Your implementation of the FirestoreStore struct **must** make all the tests in the provided test file pass.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer with deep experience using the Google Cloud Firestore client library, particularly with transactions and batch operations.

**Overall Goal**: We have an executable specification in the form of a failing integration test (firestore\_test.go). Your task is now to write the production-ready Go code that implements the routing.MessageStore interface and satisfies that specification.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **internal/platform/persistence/firestore.go**.

**Key Requirements to Fulfill**:

* **Pass All Tests**: This is the most critical requirement.  
* **Firestore Transactions (L2-2.1)**: Both the StoreMessages and DeleteMessages methods **must** use client.RunTransaction to ensure that the batch operations are atomic. This is non-negotiable.  
* **Collection Path (L3-2)**: The implementation **must** store and retrieve documents from the correct hierarchical path: user-messages/{recipient-urn-string}/messages/{message-id}.  
* **Error Handling**:  
  * The constructor (NewFirestoreStore) **must** return an error if the provided Firestore client is nil.  
  * RetrieveMessages **must** correctly handle cases where a user's document or messages sub-collection does not exist by returning an empty slice and a nil error.  
* **Message ID Generation**: If any SecureEnvelope passed to StoreMessages has an empty MessageID field, the implementation **must** generate and assign a new UUID to it before storing.

**Dependencies**:

1. **The Contract**: pkg/routing/interfaces\_routing.go  
2. **The Specification**: internal/platform/persistence/firestore\_test.go  
3. **The Data Models**: pkg/transport/envelope.go and pkg/urn/urn.go

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created. Mandates the use of Firestore transactions to fulfill the atomicity requirement (L2-2.1).