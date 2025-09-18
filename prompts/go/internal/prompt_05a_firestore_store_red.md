# **Prompt 05a: FirestoreStore Integration Tests (Red)**

**Objective**: Generate a complete **integration test** file (firestore\_test.go) for the FirestoreStore. This test will serve as the executable specification for a FirestoreStore implementation that does not yet exist.

**TDD Phase**: Red (Generate failing tests before the code exists).

## **Primary Prompt (Concise)**

Your task is to generate the Go code for the internal/platform/persistence/firestore\_test.go file.

This file must contain an integration test suite for a FirestoreStore struct that will implement the routing.MessageStore interface. The tests **must** use a real Firestore client connected to a **Firestore emulator** to validate the full lifecycle: StoreMessages, RetrieveMessages, and DeleteMessages.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer specializing in writing robust integration tests for cloud services using official emulators.

**Overall Goal**: We are beginning the TDD implementation phase of the routing service. Our first component is the concrete Firestore adapter that implements the MessageStore interface. We will define its behavior by writing a comprehensive integration test first.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **internal/platform/persistence/firestore\_test.go**.

**Key Requirements to Test**:

* **Emulator Setup**: The test **must** programmatically set up and connect to a Firestore emulator for the duration of the test suite.  
* **Full Lifecycle**: The primary test case **must** validate the complete, sequential user story:  
  1. **Store**: StoreMessages is called to persist a slice of two SecureEnvelopes for a specific user URN.  
  2. **Verify Storage**: Use the emulator's client directly to confirm that two documents now exist in the correct collection path (user-messages/{urn}/messages/{msgId}).  
  3. **Retrieve**: RetrieveMessages is called for that same user, and the test **must assert** that the two retrieved envelopes are identical to the ones that were stored.  
  4. **Delete**: DeleteMessages is called with the IDs of the retrieved messages.  
  5. **Verify Deletion**: A final call to RetrieveMessages **must** return an empty slice and a nil error.  
* **Edge Cases**: You **must** also include tests for:  
  * Retrieving messages for a user who has never had any messages stored. This **must** succeed and return an empty slice.  
  * Storing an empty slice of messages for a user. This **must** succeed and be a no-op.

**Dependencies**:

* You should assume the pkg/routing, pkg/transport, and pkg/urn packages, containing the interfaces and data models, are available for import.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created. Focuses on a full-lifecycle integration test using a Firestore emulator as the primary specification.

---